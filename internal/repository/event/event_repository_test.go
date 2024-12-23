package event_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository/event"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := event.NewEventRepository(gormDB, logger)

	expectedTime := time.Date(0, 1, 1, 15, 4, 5, 0, time.UTC)
	expectedEvent := &entity.Event{
		ID:          1,
		Name:        "Sample Event",
		Description: "This is a test event.",
		Date:        time.Now().UTC(),
		Time:        helper.SQLTime(expectedTime),
		VenueID:     2,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `events` (`name`,`description`,`date`,`time`,`venue_id`,`created_at`,`updated_at`,`deleted_at`,`id`) VALUES (?,?,?,?,?,?,?,?,?)")).
		WithArgs(expectedEvent.Name, expectedEvent.Description, expectedEvent.Date, expectedEvent.Time, expectedEvent.VenueID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)) 
	mock.ExpectCommit()

	err = repo.Create(gormDB, expectedEvent)

	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}


func TestRepositoryImpl_Update(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := event.NewEventRepository(gormDB, logger)

	expectedTime := time.Date(0, 1, 1, 15, 4, 5, 0, time.UTC)
	expectedEvent := &entity.Event{
		ID:          1,
		Name:        "Sample Event",
		Description: "This is a test event.",
		Date:        time.Now().UTC(),
		Time:        helper.SQLTime(expectedTime),
		VenueID:     2,
	}

	// Mock the query for Update
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `events` SET `name`=?,`description`=?,`date`=?,`time`=?,`venue_id`=?,`created_at`=?,`updated_at`=?,`deleted_at`=? WHERE `events`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs(expectedEvent.Name, expectedEvent.Description, expectedEvent.Date, expectedEvent.Time, expectedEvent.VenueID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, expectedEvent.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	err = repo.Update(gormDB, expectedEvent)

	// Assertions
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}


func TestRepositoryImpl_Delete(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := event.NewEventRepository(gormDB, logger)

	// Define expected event
	expectedEvent := &entity.Event{
		ID: 1,
	}

	// Mock the query for Delete
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta( "UPDATE `events` SET `deleted_at`=? WHERE `events`.`id` = ? AND `events`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), expectedEvent.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	err = repo.Delete(gormDB, expectedEvent)

	// Assertions
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestEventRepository_GetByID(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()

	repo := event.NewEventRepository(gormDB, logger)
	expectedTime := time.Date(0, 1, 1, 15, 4, 5, 0, time.UTC)
	expectedEvent := &entity.Event{
		ID:          1,
		Name:        "Sample Event",
		Description: "This is a test event.",
		Date:        time.Now().UTC(),
		Time:        helper.SQLTime(expectedTime),
		VenueID:     2,
	}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "date", "time", "venue_id"}).
		AddRow(expectedEvent.ID, expectedEvent.Name, expectedEvent.Description, expectedEvent.Date, expectedEvent.Time, expectedEvent.VenueID)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `events` WHERE id = ? AND `events`.`deleted_at` IS NULL ORDER BY `events`.`id` LIMIT ?")).
		WithArgs(expectedEvent.ID, 1).
		WillReturnRows(rows)

	data := &entity.Event{}
	err = repo.GetByID(gormDB, data, expectedEvent.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvent.ID, data.ID)
	assert.Equal(t, expectedEvent.Name, data.Name)
	assert.Equal(t, expectedEvent.Description, data.Description)
	assert.Equal(t, expectedEvent.Date, data.Date)
	assert.Equal(t, expectedEvent.Time, data.Time)
	assert.Equal(t, expectedEvent.VenueID, data.VenueID)
	mock.ExpectationsWereMet()
}

func TestEventRepository_GetPaginated(t *testing.T) {
	// Create SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Initialize repository
	logger := logrus.New()
	repo := event.NewEventRepository(gormDB, logger)

	// Menyiapkan query options untuk pagination
	opts := &model.EventQueryOptions{
		Page: 1,
		Size: 10,
	}

	// Define expected events
	expectedTime := time.Date(0, 1, 1, 15, 4, 5, 0, time.UTC)
	expectedEvents := []entity.Event{
		{ID: 1, Name: "Event 1", Description: "First event", Date: time.Now(), Time: helper.SQLTime(expectedTime), VenueID: 1},
		{ID: 2, Name: "Event 2", Description: "Second event", Date: time.Now(), Time: helper.SQLTime(expectedTime), VenueID: 2},
	}

	// Mock rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "date", "time", "venue_id"}).
		AddRow(expectedEvents[0].ID, expectedEvents[0].Name, expectedEvents[0].Description, expectedEvents[0].Date, expectedEvents[0].Time, expectedEvents[0].VenueID).
		AddRow(expectedEvents[1].ID, expectedEvents[1].Name, expectedEvents[1].Description, expectedEvents[1].Date, expectedEvents[1].Time, expectedEvents[1].VenueID)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `events` WHERE `events`.`deleted_at` IS NULL")).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `events` WHERE `events`.`deleted_at` IS NULL LIMIT ?")).
		WithArgs(opts.Size).
		WillReturnRows(rows)

	var data []entity.Event
	totalCount, err := repo.GetPaginated(gormDB, &data, opts)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, int64(2), totalCount)
	assert.Equal(t, expectedEvents[0].ID, data[0].ID)
	assert.Equal(t, expectedEvents[1].ID, data[1].ID)
	mock.ExpectationsWereMet()
}
