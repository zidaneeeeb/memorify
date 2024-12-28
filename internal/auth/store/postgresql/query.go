package postgresql

const (
	queryCreateUser = `
		INSERT INTO 
			user_info
		(
			id,
			fullname,
			username,
			email,
			type,
			quota,
			create_time
		)
		VALUES
			(
				:id,
				:fullname,
				:username,
				:email,
				:type,
				:quota,
				:create_time
			)
		RETURNING
			id
	`

	queryGetUser = `
		SELECT
			id,
			fullname,
			username,
			email,
			type,
			quota,
			create_time,
			update_time
		FROM
			user_info
		%s
	`

	queryUpdateUser = `
		UPDATE
			user_info
		SET
			fullname = :fullname,
			username = :username,
			email = :email,
			type = :type,
			quota = :quota,
			update_time = :update_time
		WHERE
			id = :id
	`
)
