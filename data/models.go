package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	// ID        int    `json:"id"`
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	// Password  string    `json:"-"`
	Password  string    `json:"password"`
	Active    bool      `json:"user_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Info().Msgf("Error scanning: %v", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`
	// query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = ($1)::uuid`
	// query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = UUID(?)`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Check the provided user "id" does exist or not
// func (u *User) CheckId(id string) error {
// Since we're using uuid instead of int id in db
func (u *User) CheckId(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// query := `if exists(select * from users where id = $1)`
	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetOne returns one user by id
// func (u *User) GetOne(id int) (*User, error) {
// func (u *User) CheckId(id string) error {
func (u *User) GetOne(id string) (*User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	`

	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	// if you create a unique key constraint on title & body columns, you can use insert statement as below to ignore if record
	// already exists
	// insert into posts(id, title, body) values (1, 'First post', 'Awesome') on conflict (title, body) do nothing;

	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// GetAll returns a slice of all users, sorted by last name for pagination
// It would require limit and cursor (timestamp)
func (u *User) GetAllForPagination(cursorTime time.Time, cursorUUID string, isFirstQuery bool) ([]*User, error) {
	// newTime := cursorTime
	// if cursorTime == time.Now() {
	// 	newTime = time.Date(1970, time.Month(1), 0, 0, 0, 0, 0, time.UTC) // epoch time: 1970, 1 January 00:00:00 UTC
	// } else {
	// 	newTime = cursorTime
	// }

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	log.Info().Msgf("from GetAllForPagination: cursorTime:%v", cursorTime)

	// query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	// from users where created_at <= $1 and id < $2 order by created_at desc, id desc limit $3`
	// query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	// from users where created_at > $1 and id > $2 order by created_at desc, id desc limit $3`

	var query string
	var limit, offset int

	newCursorTime := cursorTime.String()
	newCursorTimeSlice := strings.Split(newCursorTime, "+")
	newCursorTime = strings.TrimSpace(newCursorTimeSlice[0])

	log.Info().Msgf("from GetAllForPagination: newCursorTime:%v", newCursorTime)

	if isFirstQuery {
		query = `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users where created_at > $1 order by created_at asc limit $2 offset $3`
		limit = 10
		offset = 0 // we can remove the offset totally as it's not needed now
	} else {
		query = `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users where created_at > $1 order by created_at asc limit $2 offset $3`
		limit = 10
		offset = 0 // we can remove the offset totally as it's not needed now
	}

	// rows, err := db.QueryContext(ctx, query, newTime, cursorUUID, 10) // Here 10 is limit. By default limit is 10 per page if no value provided
	rows, err := db.QueryContext(ctx, query, newCursorTime, limit, offset) // Here 10 is limit. By default limit is 10 per page if no value provided
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Info().Msgf("Error scanning: %v", err)
			return nil, err
		}

		// Don't want to return the user's password so setting the value to "-"
		user.Password = "-"
		users = append(users, &user)
	}

	return users, nil
}
