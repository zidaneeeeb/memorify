package postgresql

const (
	queryCreateContent = `
		INSERT INTO
			content
			(
				user_id,
				template_id,
				detail_content_json_text,
				status,
				create_time
			)
		VALUES
			(
				:user_id,
				:template_id,
				:detail_content_json_text,
				:status,
				:create_time
			)
		RETURNING
			id
	`

	queryGetContent = `
		SELECT
			c.id,
			c.user_id,
			u.fullname as user_name,
			c.template_id,
			t.label as template_label,
			t.name as template_name,
			c.detail_content_json_text,
			c.status,
			c.create_time,
			c.update_time
		FROM
			content c
		LEFT JOIN
			template t
		ON 
			t.id = c.template_id
		Left JOIN
			user_info u
		ON
			u.id = c.user_id
		%s
	`

	queryUpdateContent = `
		UPDATE
			content
		SET
			template_id = :template_id,
			detail_content_json_text = :detail_content_json_text,
			status = :status,
			update_time = :update_time
		WHERE
			id = :id 
	`

	queryDeleteContent = `
		DELETE FROM
			content
		WHERE
			id = :id
		AND 
			user_id = :user_id
	`
)
