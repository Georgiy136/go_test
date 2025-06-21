CREATE OR REPLACE FUNCTION goods_upd(_src JSONB) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _dt          TIMESTAMPTZ := now() AT TIME ZONE 'Europe/Moscow';
    _good_id     INTEGER;
    _project_id  INTEGER;
    _name        VARCHAR(32);
    _description TEXT;
    _priority    INTEGER;
    _deleted_at  TIMESTAMPTZ;
    _js          JSONB;
BEGIN
SELECT s.good_id,
       s.project_id,
       s.name,
       s.description,
       s.priority,
       s.deleted_at
INTO _good_id,
     _project_id,
     _name,
     _description,
     _priority,
     _deleted_at
FROM jsonb_to_record(_src) as s (good_id INTEGER,
                                 project_id INTEGER,
                                 name VARCHAR(32),
                                 description TEXT,
                                 priority INTEGER,
                                 deleted_at TIMESTAMPTZ
    )
         LEFT JOIN goods g ON s.good_id = g.good_id;

IF _good_id IS NOT NULL AND _project_id IS NOT NULL THEN
        IF NOT EXISTS (SELECT 1 FROM goods g WHERE g.good_id = _good_id AND g.project_id = _project_id) THEN
            RAISE EXCEPTION 'goods not found. good_id = %, project_id = %', _good_id, _project_id;
        END IF;
END IF;

    IF _good_id IS NULL THEN _good_id = nextval('good_sq');
END IF;

WITH ins_cte AS (
INSERT INTO goods AS g (good_id,
                        project_id,
                        name,
                        description,
                        priority,
                        created_at,
                        deleted_at)
SELECT _good_id,
       _project_id,
       _name,
       _description,
       _priority,
       _dt,
       _deleted_at
    ON CONFLICT (good_id, project_id) DO UPDATE
                                      SET name        = excluded.name,
                                          description = excluded.description,
                                          priority    = excluded.priority,
                                          created_at  = excluded.created_at,
                                          deleted_at  = excluded.deleted_at
                                      RETURNING g.*)

SELECT jsonb_build_object('data', row_to_json(ins_cte))
FROM ins_cte
    INTO _js;
RETURN _js;
END
$$;
