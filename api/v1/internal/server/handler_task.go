package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/osetr/rest-api/api/v1/internal/model"
)

// Add new task
func (s *APIServer) createTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Task{}
		var err error

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if err = req.Validate(); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}

		var position int
		t := &model.Task{
			ID:          genID(taskID),
			ColumnID:    req.ColumnID,
			Name:        req.Name,
			Description: req.Description,
			Position:    &position,
		}

		if err = s.Store.DB.QueryRow(
			"SELECT coalesce(MAX(position), -1)+1 FROM tasks WHERE column_id = $1",
			req.ColumnID,
		).Scan(&position); err != nil {
			s.errorresp(w, r, http.StatusNotFound, err)
			return
		}

		res, err := s.Store.DB.Exec(
			"INSERT INTO tasks (id, column_id, name, description, position) VALUES ($1, $2, $3, $4, $5)",
			t.ID,
			t.ColumnID,
			t.Name,
			t.Description,
			t.Position,
		)
		if err != nil {
			s.errorresp(w, r, http.StatusUnprocessableEntity, errors.New("column not found"))
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusInternalServerError, errors.New("task not created"))
			return
		}

		s.successresp(w, r, http.StatusCreated, t)
	}
}

// Get tasks list
func (s *APIServer) getTasks() http.HandlerFunc {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query()
		pageSize := int(10)
		page := int(0)
		var err error

		if _, present := query["page_size"]; present {
			ps, err := strconv.ParseInt(query.Get("page_size"), 10, 64)
			if err != nil || ps < 0 {
				s.errorresp(w, r, http.StatusBadRequest, errors.New("uncorrect page_size param"))
				return
			}
			pageSize = int(ps)
		}
		if _, present := query["page"]; present {
			p, err := strconv.ParseInt(query.Get("page"), 10, 64)
			if err != nil || p < 0 {
				s.errorresp(w, r, http.StatusBadRequest, errors.New("uncorrect paga param"))
				return
			}
			page = int(p)
		}

		var rows *sql.Rows
		if _, present := query["column_id"]; present {
			rows, err = s.Store.DB.Query(
				"SELECT * FROM tasks WHERE column_id = $1 ORDER BY position",
				query.Get("column_id"),
			)
		} else {
			rows, err = s.Store.DB.Query("SELECT * FROM tasks")
		}

		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		result := []model.Task{}
		for rows.Next() {
			c := &model.Task{}
			err = rows.Scan(
				&c.ID,
				&c.ColumnID,
				&c.Name,
				&c.Description,
				&c.Position)
			result = append(result, *c)
		}

		pstart := min(page*pageSize, int(len(result)))
		pfinish := min(page*pageSize+pageSize, int(len(result)))
		s.successresp(w, r, http.StatusOK, map[string]interface{}{
			"total_count": len(result),
			"data":        result[pstart:pfinish],
		})
	}
}

// Get single task
func (s *APIServer) getTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		t := &model.Task{}
		if err := s.Store.DB.QueryRow(
			"SELECT id, column_id, name, description, position FROM tasks WHERE id = $1",
			params["id"],
		).Scan(
			&t.ID,
			&t.ColumnID,
			&t.Name,
			&t.Description,
			&t.Position,
		); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("task not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, t)
	}
}

// Patch task
func (s *APIServer) patchTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Task{}
		params := mux.Vars(r)
		var err error

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if err = req.ValidatePatch(); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}

		id := params["id"]
		t := &model.Task{
			ID:          &id,
			ColumnID:    req.ColumnID,
			Name:        req.Name,
			Description: req.Description,
			Position:    req.Position,
		}

		var res sql.Result
		if t.Description == nil && t.Name != nil {
			res, err = s.Store.DB.Exec(
				"UPDATE tasks SET name = $2 WHERE id = $1",
				t.ID,
				t.Name,
			)
		} else if t.Description != nil && t.Name == nil {
			res, err = s.Store.DB.Exec(
				"UPDATE tasks SET description = $2 WHERE id = $1",
				t.ID,
				t.Description,
			)
		} else if t.Description != nil && t.Name != nil {
			res, err = s.Store.DB.Exec(
				"UPDATE tasks SET name = $2, description = $3 WHERE id = $1",
				t.ID,
				t.Name,
				t.Description,
			)
		} else {
			err = errors.New("at least one field have to be not blank")
		}
		if err != nil {
			s.errorresp(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("task not found"))
			return
		}

		if err = s.Store.DB.QueryRow(
			"SELECT id, column_id, name, description, position FROM tasks WHERE id = $1",
			params["id"],
		).Scan(
			&t.ID,
			&t.ColumnID,
			&t.Name,
			&t.Description,
			&t.Position,
		); err != nil {
			s.errorresp(w, r, http.StatusNotFound, err)
			return
		}

		s.successresp(w, r, http.StatusOK, t)
	}
}

// Delete task
func (s *APIServer) deleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tx, err := s.Store.DB.Begin()
		params := mux.Vars(r)
		var column string
		var position int

		if err := tx.QueryRow(
			"DELETE FROM tasks WHERE id = $1 RETURNING column_id, position",
			params["id"],
		).Scan(&column, &position); err != nil {
			s.errorresp(w, r, http.StatusUnprocessableEntity, errors.New("task not found"))
			return
		}

		if _, err = tx.Exec(
			"UPDATE tasks SET position = position - 1 WHERE column_id = $1 AND position > $2",
			column,
			position,
		); err != nil {
			s.errorresp(w, r, http.StatusUnprocessableEntity, err)
			tx.Rollback()
			return
		}

		err = tx.Commit()
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		s.successresp(w, r, http.StatusNoContent, nil)
	}
}

// Up task
func (s *APIServer) upTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var position int

		if err := s.Store.DB.QueryRow(
			`WITH
				curtask AS 
					(SELECT * FROM tasks WHERE id = $1)

			UPDATE tasks
			SET position = 
				CASE
					WHEN position = (SELECT position FROM curtask) - 1 THEN position + 1
					WHEN id = $1 AND position != 0 THEN position - 1
					ELSE position
				END
			WHERE
				(
					column_id = (SELECT column_id FROM curtask)
				AND
					(
						position = (SELECT position FROM curtask) - 1
					OR
						position = (SELECT position FROM curtask)
					)
				)
			RETURNING (SELECT position FROM curtask)`,
			params["id"],
		).Scan(&position); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("task not found"))
			return
		}
		if position == 0 {
			s.errorresp(w, r, http.StatusBadRequest, errors.New("can't move up"))
			return
		}

		s.successresp(w, r, http.StatusOK, map[string]interface{}{"detail": "successfully moved"})
	}
}

// Down task
func (s *APIServer) downTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var maxPosition int
		var position int

		if err := s.Store.DB.QueryRow(
			`WITH
				curtask AS 
					(SELECT * FROM tasks WHERE id = $1),
				curtasks AS 
					(
						SELECT * FROM tasks WHERE
							column_id = (SELECT column_id FROM curtask)
					)

			UPDATE tasks
			SET position = 
				CASE
					WHEN position = (SELECT position FROM curtask) + 1 THEN position - 1
					WHEN 
						id = $1
					AND
						position !=
							(SELECT MAX(position) FROM curtasks)
					THEN position + 1
					ELSE position
				END
			WHERE
				(
					column_id = (SELECT column_id FROM curtask)
				AND
					(
						position = (SELECT position FROM curtask) + 1
					OR
						position = (SELECT position FROM curtask)
					)
				)
			RETURNING (SELECT position FROM curtask), (SELECT MAX(position) FROM curtasks)`,
			params["id"],
		).Scan(&position, &maxPosition); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("task not found"))
			return
		}
		if position == maxPosition {
			s.errorresp(w, r, http.StatusBadRequest, errors.New("can't move down"))
			return
		}

		s.successresp(w, r, http.StatusOK, map[string]interface{}{"detail": "successfully moved"})
	}
}

// Left task
func (s *APIServer) leftTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		res, err := s.Store.DB.Exec(
			`WITH
				curtask AS 
					(SELECT * from tasks WHERE id = $1),
				curcol AS
					(SELECT * FROM columns WHERE id = 
						(SELECT column_id FROM curtask)
					),
				goalcol AS
					(SELECT * FROM columns WHERE
						project_id =
							(SELECT project_id FROM curtask)
					AND
						position=
							(SELECT position-1 FROM curcol)
					)

			UPDATE tasks
			SET position =
					CASE
						WHEN position > (SELECT position FROM curtask) THEN position-1
						WHEN position = (SELECT position FROM curtask) THEN
							(SELECT coalesce(MAX(position), -1)+1 FROM tasks WHERE column_id = 
									(SELECT id FROM goalcol)
							)
						ELSE position
					END,
				column_id = 
					CASE
						WHEN id = $1 THEN (SELECT id FROM goalcol)
						ELSE column_id
					END
			WHERE column_id = (SELECT column_id FROM curtask)`,
			params["id"],
		)
		if err != nil {
			s.errorresp(w, r, http.StatusBadRequest, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("task not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, map[string]interface{}{"detail": "successfully moved"})
	}
}

// Right task
func (s *APIServer) rightTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		res, err := s.Store.DB.Exec(
			`WITH
				curtask AS 
					(SELECT * from tasks WHERE id = $1),
				curcol AS
					(SELECT * FROM columns WHERE id = 
						(SELECT column_id FROM curtask)
					),
				goalcol AS
					(SELECT * FROM columns WHERE
						project_id =
							(SELECT project_id FROM curtask)
					AND
						position=
							(SELECT position+1 FROM curcol)
					)

			UPDATE tasks
			SET position =
					CASE
						WHEN position > (SELECT position FROM curtask) THEN position-1
						WHEN position = (SELECT position FROM curtask) THEN
							(SELECT coalesce(MAX(position), -1)+1 FROM tasks WHERE column_id = 
									(SELECT id FROM goalcol)
							)
						ELSE position
					END,
				column_id = 
					CASE
						WHEN id = $1 THEN (SELECT id FROM goalcol)
						ELSE column_id
					END
			WHERE column_id = (SELECT column_id FROM curtask)`,
			params["id"],
		)
		if err != nil {
			s.errorresp(w, r, http.StatusBadRequest, errors.New("can't right task"))
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(w, r, http.StatusNotFound, errors.New("task not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, map[string]interface{}{"detail": "successfully moved"})
	}
}
