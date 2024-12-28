package postgresql

const (
	queryCreatePayment = `
		INSERT INTO
			payment
			(
				user_id,
				content_id,
				amount,
				proof_payment_url,
				date,
				status,
				create_time
			)
		VALUES
			(
				:user_id,
				:content_id,
				:amount,
				:proof_payment_url,
				:date,
				:status,
				:create_time
			)
		RETURNING
			id
	`

	queryGetPayment = `
		SELECT
			p.id,
			p.user_id,
			u.fullname as user_name,
			u.type as user_type,
			u.quota as user_quota,
			p.content_id,
			t.label as template_label,
			t.name as template_name,
			p.amount,
			p.proof_payment_url,
			p.date,
			p.status,
			p.create_time,
			p.update_time
		FROM
			payment p
		LEFT JOIN
			user_info u
		ON
			u.id = p.user_id
		LEFT JOIN
			content c
		ON
			c.id = p.content_id
		LEFT JOIN
			template t
		ON
			t.id = c.template_id
		%s
	`

	queryUpdatePayment = `
		UPDATE
			payment
		SET
			user_id = :user_id,
			content_id = :content_id,
			amount = :amount,
			proof_payment_url = :proof_payment_url,
			date = :date,
			status = :status,
			update_time = :update_time
		WHERE
			id = :id
	`
)
