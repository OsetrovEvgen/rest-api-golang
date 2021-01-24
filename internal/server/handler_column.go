package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/osetr/rest-api/internal/model"
)

// Add new column
func (s *APIServer) createColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Column{}
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
		c := &model.Column{
			ID:        genID(columnID),
			ProjectID: req.ProjectID,
			Name:      req.Name,
			Position:  &position,
		}

		if err = s.Store.DB.QueryRow(
			"SELECT MAX(position) FROM columns WHERE project_id = $1",
			req.ProjectID,
		).Scan(&position); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("project not found"))
			return
		}
		position++

		if _, err = s.Store.DB.Exec(
			"INSERT INTO columns (id, project_id, name, position) VALUES ($1, $2, $3, $4)",
			c.ID,
			c.ProjectID,
			c.Name,
			c.Position,
		); err != nil {
			s.errorresp(w, r, http.StatusBadRequest, errors.New("name in project must be unique"))
			return
		}

		s.successresp(w, r, http.StatusCreated, c)
	}
}

// Get columns list
func (s *APIServer) getColumns() http.HandlerFunc {
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
		if _, present := query["project_id"]; present {
			rows, err = s.Store.DB.Query(
				"SELECT * FROM columns WHERE project_id = $1 ORDER BY position",
				query.Get("project_id"),
			)
		} else {
			rows, err = s.Store.DB.Query("SELECT * FROM columns")
		}

		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		result := []model.Column{}
		for rows.Next() {
			c := &model.Column{}
			err = rows.Scan(
				&c.ID,
				&c.ProjectID,
				&c.Name,
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

// Get single column
func (s *APIServer) getColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		c := &model.Column{}
		if err := s.Store.DB.QueryRow(
			"SELECT id, project_id, name, position FROM columns WHERE id = $1",
			params["id"],
		).Scan(
			&c.ID,
			&c.ProjectID,
			&c.Name,
			&c.Position,
		); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("column not found"))
			return
		}

		s.successresp(w, r, http.StatusOK, c)
	}
}

// Patch column
func (s *APIServer) patchColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &model.Column{}
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
		c := &model.Column{
			ID:   &id,
			Name: req.Name,
		}

		var res sql.Result
		if c.Name != nil {
			res, err = s.Store.DB.Exec(
				"UPDATE columns SET name = $2 WHERE id = $1",
				c.ID,
				c.Name,
			)
		} else {
			s.errorresp(w, r,
				http.StatusBadRequest,
				errors.New("name field have to be not blank"),
			)
			return
		}
		if err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}
		if st, _ := res.RowsAffected(); st == 0 {
			s.errorresp(
				w, r,
				http.StatusNotFound,
				errors.New("column not found"),
			)
			return
		}

		if err = s.Store.DB.QueryRow(
			"SELECT id, project_id, name, position FROM columns WHERE id = $1",
			params["id"],
		).Scan(
			&c.ID,
			&c.ProjectID,
			&c.Name,
			&c.Position,
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}

		s.successresp(w, r, http.StatusOK, c)
	}
}

// Delete column
func (s *APIServer) deleteColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tx, err := s.Store.DB.Begin()
		params := mux.Vars(r)
		var count int

		if err = tx.QueryRow(
			`SELECT COUNT (*) FROM columns WHERE 
				project_id = 
					(SELECT project_id FROM columns WHERE id = $1)`,
			params["id"],
		).Scan(&count); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("column not found"))
			return
		}
		if count == 1 {
			s.errorresp(
				w, r,
				http.StatusBadRequest,
				errors.New("can't delete the only column in project"),
			)
			return
		}

		if _, err = tx.Exec(
			`WITH
				curcol AS
					(SELECT * FROM columns WHERE id = $1 ),
				goalcol AS
					(
						SELECT * FROM columns WHERE
							project_id = (SELECT project_id FROM curcol)
						AND
							position =
								CASE
									WHEN  (SELECT position FROM curcol) = 0 THEN 1
									ELSE (SELECT position FROM curcol) - 1
								END
					)

			UPDATE tasks
			SET
				column_id = (SELECT id FROM goalcol),
				position = position +
					(
						SELECT coalesce(MAX(position), -1)+1 FROM tasks WHERE
							column_id = (SELECT id FROM goalcol)
					)
			WHERE column_id = $1`,
			params["id"],
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
			return
		}

		var position int
		var project string
		if err = tx.QueryRow(
			"DELETE FROM columns WHERE id = $1 RETURNING project_id, position",
			params["id"],
		).Scan(&project, &position); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("column not found"))
			tx.Rollback()
			return
		}

		if _, err = tx.Exec(
			"UPDATE columns SET position = position - 1 WHERE project_id = $1 AND position > $2",
			project,
			position,
		); err != nil {
			s.errorresp(w, r, http.StatusInternalServerError, err)
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

// Left column
func (s *APIServer) leftColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var position int

		if err := s.Store.DB.QueryRow(
			`WITH
				curcol AS 
					(SELECT * FROM columns WHERE id = $1)

			UPDATE columns
			SET position = 
					CASE
						WHEN position = (SELECT position FROM curcol) - 1 THEN position + 1
						WHEN id = $1 AND position != 0 THEN position - 1
						ELSE position
					END
			WHERE
				(
					project_id = (SELECT project_id FROM curcol)
				AND
					(
						position = (SELECT position FROM curcol) - 1
					OR
						position = (SELECT position FROM curcol)
					)
				)
			RETURNING (SELECT position FROM curcol)`,
			params["id"],
		).Scan(&position); err != nil {
			s.errorresp(w, r, http.StatusNotFound, errors.New("column not found"))
			return
		}
		if position == 0 {
			s.errorresp(w, r, http.StatusBadRequest, errors.New("can't move left"))
			return
		}

		s.successresp(w, r, http.StatusOK, map[string]interface{}{"detail": "successfully moved"})
	}
}

// Right column
func (s *APIServer) rightColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var maxPosition int
		var position int

		if err := s.Store.DB.QueryRow(
			`WITH
				curcol AS 
					(SELECT * FROM columns WHERE id = $1),
				curcols AS
					(
					SELECT * FROM columns WHERE
						project_id = (SELECT project_id FROM curcol)
					)

			UPDATE columns
			SET position = 
					CASE
						WHEN position = (SELECT position FROM curcol) + 1 THEN position - 1
						WHEN 
							(id = $1
						AND
							position !=
								(SELECT MAX(position) FROM curcols))
						THEN position + 1
						ELSE position
					END
			WHERE
				(
					project_id = (SELECT project_id FROM curcol)
				AND
					(
						position = (SELECT position FROM curcol) + 1
					OR
						position = (SELECT position FROM curcol)
					)
				)
			RETURNING (SELECT position FROM curcol), (SELECT MAX(position) FROM curcols)`,
			params["id"],
		).Scan(&position, &maxPosition); err != nil {
			s.errorresp(w, r, http.StatusNotFound, err)
			return
		}
		if position == maxPosition {
			s.errorresp(w, r, http.StatusBadRequest, errors.New("can't move right"))
			return
		}

		s.successresp(w, r, http.StatusOK, map[string]interface{}{"detail": "successfully moved"})
	}
}
