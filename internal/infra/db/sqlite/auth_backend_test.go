package sqlite

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

func TestAuthBackendSave(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	uow := createUnitOfWork(db)
	backend := NewSqliteAuthBackend(db)
	user := createUser(uuid.New(), "jeff@amazon.com", time.Now(), "", time.Unix(0, 0))

	ctx, cancel := createContextWith5MinutesTimeout()
	defer cancel()

	if err := backend.Save(ctx, uow, user); err != nil {
		t.Fatalf("Save(): expected ok got err %v", err)
	}
}

func TestAuthBackendSaveDuplicate(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	uow := createUnitOfWork(db)
	backend := NewSqliteAuthBackend(db)
	user := createUser(uuid.New(), "jeff@amazon.com", time.Now(), "", time.Unix(0, 0))

	ctx, cancel := createContextWith5MinutesTimeout()
	defer cancel()

	if err := backend.Save(ctx, uow, user); err != nil {
		t.Fatalf("Save(): expected ok got err %v", err)
	}

	err := backend.Save(ctx, uow, user)
	if err == nil {
		t.Fatalf("Save(): expected duplicate err got ok")
	}
	if !errors.Is(err, auth.ErrUserAlreadyExists) {
		t.Fatalf("Save(): expected %s got %s", auth.ErrUserAlreadyExists, err)
	}
}

func TestAuthBackendSaveDifferentIdSameEmail(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	uow := createUnitOfWork(db)
	backend := NewSqliteAuthBackend(db)
	user := createUser(uuid.New(), "jeff@amazon.com", time.Now(), "", time.Unix(0, 0))

	ctx, cancel := createContextWith5MinutesTimeout()
	defer cancel()

	if err := backend.Save(ctx, uow, user); err != nil {
		t.Fatalf("Save(): expected ok got err %v", err)
	}

	anotherUser := createUser(uuid.New(), "jeff@amazon.com", time.Now(), "", time.Unix(0, 0))
	err := backend.Save(ctx, uow, anotherUser)
	if err == nil {
		t.Fatalf("Save(): expected duplicate err got ok")
	}
	if !errors.Is(err, auth.ErrUserAlreadyExists) {
		t.Fatalf("Save(): expected %s got %s", auth.ErrUserAlreadyExists, err)
	}

}

func TestAuthBackendGetByEmailOrSessionId(t *testing.T) {
	testcases := []struct {
		user          auth.User
		findByEmail   bool
		email         string
		sessionId     string
		shouldBeFound bool
	}{
		{
			user:          createUser(uuid.New(), "bill@microsoft.com", time.Now(), "", time.Unix(0, 0)),
			findByEmail:   true,
			email:         "bill@microsoft.com",
			sessionId:     "",
			shouldBeFound: true,
		},
	}

	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	uow := createUnitOfWork(db)
	backend := NewSqliteAuthBackend(db)

	for _, test := range testcases {
		t.Run(test.user.Email.Value(), func(t *testing.T) {
			ctx, cancel := createContextWith5MinutesTimeout()
			defer cancel()

			// pre-insert data
			if test.shouldBeFound {
				if err := backend.Save(ctx, uow, test.user); err != nil {
					t.Fatalf("Save(): expected ok got err %v", err)
				}
			}

			var user *auth.User
			var found bool
			var errGet error

			if test.findByEmail {
				email, err := valueobjects.NewEmail(test.email)
				if err != nil {
					t.Fatalf("could not create email (%s): %s", test.email, err)
				}
				user, found, errGet = backend.GetByEmail(ctx, uow, email)
			} else {
				user, found, errGet = backend.GetBySessionId(ctx, uow, test.sessionId)
			}

			if errGet != nil {
				if !test.shouldBeFound && errors.Is(errGet, auth.ErrUserNotFound) {
					// success
					return
				}
				t.Fatalf("GetBy(Email|SessionId)(): expected ok got err %v", errGet)
			}
			if !found && !test.shouldBeFound {
				// succeess
				return
			}
			if !found && test.shouldBeFound {
				t.Fatal("GetBy(Email|SessionId)(): expected to find user, found nothing")
			}
			if !reflect.DeepEqual(*user, test.user) {
				t.Fatalf("GetBy(Email|SessionId)(): expected %v got %v", test.user, *user)
			}
		})
	}
}

func TestAuthBackendUpdate(t *testing.T) {
	testcases := []struct {
		title                   string
		user                    auth.User
		updateEmail             string     // empty if no update
		updateSessionId         *string    // nil if no update
		updateHashedPassword    string     // empty if no update
		updateSessionExpireDate *time.Time // nil if no update
	}{
		{
			title:                   "update email",
			user:                    createUser(uuid.New(), "jeff@google.com", time.Now(), "", time.Unix(0, 0)),
			updateEmail:             "jeff@amazon.com",
			updateSessionId:         nil,
			updateHashedPassword:    "",
			updateSessionExpireDate: nil,
		},
		{
			title:                   "update session fields",
			user:                    createUser(uuid.New(), "jeff@google.com", time.Now(), "", time.Unix(0, 0)),
			updateEmail:             "",
			updateSessionId:         newStringPointer(uuid.NewString()),
			updateHashedPassword:    "",
			updateSessionExpireDate: newTimePointer(time.Now()),
		},
		{
			title:                   "update hashed password",
			user:                    createUser(uuid.New(), "jeff@google.com", time.Now(), "", time.Unix(0, 0)),
			updateEmail:             "",
			updateSessionId:         nil,
			updateHashedPassword:    "thisIsANewHashPasswordTrustMe",
			updateSessionExpireDate: nil,
		},
		{
			title:                   "update all possible fields",
			user:                    createUser(uuid.New(), "jeff@google.com", time.Now(), "", time.Unix(0, 0)),
			updateEmail:             "bill.gates@microsoft.com",
			updateSessionId:         newStringPointer(uuid.NewString()),
			updateHashedPassword:    "YetAnotherCorrectlyHashedPasswordTrustMe;-)",
			updateSessionExpireDate: newTimePointer(time.Now()),
		},
	}

	for _, test := range testcases {
		t.Run(test.title, func(t *testing.T) {
			// create a new db for each test to avoid any conflict
			db := createTmpDbWithAllMigrationsApplied()
			defer db.Close()
			uow := createUnitOfWork(db)
			backend := NewSqliteAuthBackend(db)

			ctx, cancel := createContextWith5MinutesTimeout()
			defer cancel()

			// pre-insert data
			if err := backend.Save(ctx, uow, test.user); err != nil {
				t.Fatalf("could not pre-insert data: %v", err)
			}

			// update fields
			email := test.user.Email
			if test.updateEmail != "" {
				newEmail, err := valueobjects.NewEmail(test.updateEmail)
				if err != nil {
					t.Fatalf("could not create email (%s): %s", test.updateEmail, err)
				}
				email = newEmail
			}

			sessionId := test.user.SessionId
			if test.updateSessionId != nil {
				sessionId = *test.updateSessionId
			}

			hashedPassword := test.user.HashedPassord
			if test.updateHashedPassword != "" {
				hashedPassword = auth.HashedPassword(test.updateHashedPassword)
			}

			sessionExpireDate := test.user.SessionExpireDate
			if test.updateSessionExpireDate != nil {
				sessionExpireDate = *test.updateSessionExpireDate
			}

			newUser := auth.NewUser(
				test.user.ID,
				email,
				hashedPassword,
				test.user.CreatedAt,
				sessionId,
				sessionExpireDate,
			)

			if err := backend.Update(ctx, uow, test.user, *newUser); err != nil {
				t.Fatalf("Update(): expected ok got err %v", err)
			}

			fetchedNew, found, err := backend.GetByEmail(ctx, uow, email)
			if err != nil {
				t.Fatalf("GetByEmail(): expected ok got err %v", err)
			}
			if !found {
				t.Fatalf("GetByEmail(): expected to find updated user got nothing")
			}

			if !reflect.DeepEqual(*fetchedNew, *newUser) {
				t.Fatalf("GetByEmail(): expected %v got %v", *newUser, *fetchedNew)
			}
		})
	}
}

func TestAuthBackendCantUpdateUserID(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	uow := createUnitOfWork(db)
	backend := NewSqliteAuthBackend(db)

	ctx, cancel := createContextWith5MinutesTimeout()
	defer cancel()

	user := createUser(uuid.New(), "jeff@google.com", time.Now(), "", time.Unix(0, 0))
	if err := backend.Save(ctx, uow, user); err != nil {
		t.Fatalf("Save(): expected ok got err %v", err)
	}
	newUser := auth.NewUser(
		uuid.New(),
		user.Email,
		user.HashedPassord,
		user.CreatedAt,
		user.SessionId,
		user.SessionExpireDate)

	err := backend.Update(ctx, uow, user, *newUser)
	if err == nil {
		t.Fatalf("Update(): expected err but got ok")
	}
	if !strings.Contains(err.Error(), "user ID can't be updated") {
		t.Fatalf("Update(): expected err 'user ID can't be updated' got %v", err)
	}
}

func TestAuthBackendCantUpdateToAnAlreadyTakenEmail(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	uow := createUnitOfWork(db)
	backend := NewSqliteAuthBackend(db)

	ctx, cancel := createContextWith5MinutesTimeout()
	defer cancel()

	jeff := createUser(uuid.New(), "jeff@google.com", time.Now(), "", time.Unix(0, 0))
	if err := backend.Save(ctx, uow, jeff); err != nil {
		t.Fatalf("Save(jeff): expected ok got err %v", err)
	}

	bill := createUser(uuid.New(), "bill@microsoft.com", time.Now(), "", time.Unix(0, 0))
	if err := backend.Save(ctx, uow, bill); err != nil {
		t.Fatalf("Save(bill): expected ok got err %v", err)
	}

	newBillEmail, err := valueobjects.NewEmail("jeff@google.com") // already taken by jeff
	if err != nil {
		t.Fatalf("new bill email address is invalid: %v", err)
	}
	newBill := auth.NewUser(
		bill.ID,
		newBillEmail,
		bill.HashedPassord,
		bill.CreatedAt,
		bill.SessionId,
		bill.SessionExpireDate,
	)
	err = backend.Update(ctx, uow, bill, *newBill)
	if err == nil {
		t.Fatal("Update(): expected to return err because email address is already taken")
	}
}
