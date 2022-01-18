package repository

import (
	"UserMicro/proto"
	"context"
	"database/sql"
)

type UserRepositoryServer interface {
	Create(ctx context.Context, user *proto.User) (*proto.User, error)
	ReadByID(ctx context.Context, id int64) (*proto.User, error)
	ReadByEmail(ctx context.Context, email string) (*proto.User, error)
	UpdateRole(ctx context.Context, user *proto.User, role *proto.Role) error
	GetRoleByUser(ctx context.Context, user *proto.User) (*proto.Role, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Create(ctx context.Context, user *proto.User) (*proto.User, error) {
	query := `INSERT INTO users(login_email, user_name, user_surname, role_id, password_hash)
	VALUES($1, $2, $3, $4)
	RETURNING id;`
	role, err := getRoleByPermissions(repo, map[string]bool{"is_customer": true})
	if err != nil {
		return nil, err
	}
	row := repo.db.QueryRowContext(ctx, query, user.Email, user.Name, user.Surname, role.Id, user.Password)
	err = row.Scan(&user.Id)

	return user, err
}

func (repo *UserRepo) readByField(ctx context.Context, fieldName string, value interface{}) (*proto.User, error) {
	query := `SELECT 
		usr.id, 
      	usr.login_email, 
       	usr.user_name, 
        usr.user_surname, 
        usr.password_hash, 
        rls.id as role_id,
       	rls.name as role_name,
        rls.is_admin as role_admin,
        rls.is_user as role_customer,
        rls.is_supplier as role_supplier
		FROM users as usr 
		    LEFT JOIN roles as rls 
		        ON usr.role_id = rls.id`
	query += ` WHERE usr.` + fieldName + ` = $1;`
	row := repo.db.QueryRow(query, value)
	var user proto.User
	var role proto.Role
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Surname, &user.Password,
		&role.Id, &role.Name, &role.IsAdmin, &role.IsCustomer, &role.IsSupplier)
	if err != nil {
		return &user, err
	}
	user.Role = &role

	return &user, err
}

func (repo *UserRepo) ReadByID(ctx context.Context, id int64) (*proto.User, error) {
	return repo.readByField(ctx, "id", id)
}

func (repo *UserRepo) ReadByEmail(ctx context.Context, email string) (*proto.User, error) {
	return repo.readByField(ctx, "login_email", email)
}

func (repo *UserRepo) GetRoleByUser(ctx context.Context, user *proto.User) (*proto.Role, error) {
	userDB, err := repo.readByField(ctx, "id", user.Id)
	if err != nil {
		return &proto.Role{}, err
	}
	return userDB.Role, nil
}

func (repo *UserRepo) UpdateRole(ctx context.Context, user *proto.User, role *proto.Role) error {
	permissions := map[string]bool{
		"is_admin":    role.IsAdmin,
		"is_customer": role.IsCustomer,
		"is_supplier": role.IsSupplier,
	}
	newRole, err := getRoleByPermissions(repo, permissions)
	if err != nil {
		return err
	}
	role = newRole

	query := `UPDATE users 
	SET role_id = $1
	WHERE id = $2;`
	_, err = repo.db.ExecContext(ctx, query, role.Id, user.Id)

	return err
}

func getRoleByPermissions(repo *UserRepo, permissions map[string]bool) (*proto.Role, error) {
	query := `SELECT 
	id, name
	FROM roles
	WHERE is_admin = $1 AND is_user = $2 AND is_supplier = $3`

	var role proto.Role
	for k, v := range permissions {
		switch k {
		case "is_admin":
			role.IsAdmin = v
		case "is_customer":
			role.IsCustomer = v
		case "is_supplier":
			role.IsSupplier = v
		}
	}

	row := repo.db.QueryRow(query, role.IsAdmin, role.IsCustomer, role.IsSupplier)
	err := row.Scan(&role.Id, &role.Name)

	return &role, err
}
