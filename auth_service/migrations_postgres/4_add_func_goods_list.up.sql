CREATE OR REPLACE FUNCTION goods_list(_goods_id INTEGER,
                                      _project_id INTEGER,
                                      _limit INTEGER,
                                      _offset INTEGER) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
_dt         TIMESTAMPTZ := now() AT TIME ZONE 'Europe/Moscow';
_res        JSONB;
BEGIN
        WITH goods_cte AS (SELECT g.good_id,
                                  g.project_id,
                                  g.name,
                                  p.name as project_name,
                                  g.description,
                                  g.priority,
                                  g.created_at,
                                  g.deleted_at
                           FROM goods AS g
                           INNER JOIN projects p ON p.project_id = g.project_id
                           WHERE g.good_id = COALESCE(_goods_id, g.good_id)
                             AND p.project_id = COALESCE(_project_id, g.project_id)
            LIMIT _limit OFFSET _offset)

        SELECT JSONB_BUILD_OBJECT('meta', JSONB_BUILD_OBJECT('total', (SELECT COUNT(*) FROM goods_cte),
                                                             'remove', (SELECT COUNT(*) FROM goods_cte g WHERE g.deleted_at IS NOT NULL),
                                                             'limit', _limit,
                                                             'offset', _offset),
                                  'goods', JSONB_AGG(c))
        INTO _res
        FROM (SELECT g.good_id,
                     g.project_id,
                     g.name,
                     g.project_name,
                     g.description,
                     g.priority,
                     g.created_at,
                     g.deleted_at
              FROM goods_cte g) c;

        RETURN JSONB_BUILD_OBJECT('data', _res);
END
$$;
