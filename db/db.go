package db

// TODO: THINK ABOUT SEPARATING TABLES FOR GUILD SETTINGS AND GUILD DATA/OTHER DATA

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Guild struct {
	gorm.Model
	GuildID string            `gorm:"index"`
	Members map[string]Member `gorm:"many2many:guild_members;serializer:gob"`
	GayConf GayConfig         `gorm:"type:bytes;serializer:gob"`
}

type Member struct {
	gorm.Model
	UserID    string
	Gays      uint
	SuperGays uint
}

type GayConfig struct {
	gorm.Model
	GuildID             string
	GayEmote            string
	SuperGayEmote       string
	BlocklistedChannels []string `gorm:"type:text"`
	BlocklistedUsers    []string `gorm:"type:text"`
}

type Database struct {
	db *gorm.DB
}

func (db *Database) Setup() (err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	db.db, err = gorm.Open(sqlite.Open("botkisser.db"), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return err
	}

	err = db.db.AutoMigrate(&Member{}, &Guild{})

	return err
}

func (db *Database) UpdateGuildSettings(guildID string, settings *Guild) (guild Guild, res *gorm.DB) {
	db.db.FirstOrCreate(&guild, Guild{GuildID: guildID})

	res = db.db.Model(&guild).Updates(settings)
	return
}

func (db *Database) GetGuildSettings(guildID string) (guild Guild, res *gorm.DB) {
	res = db.db.FirstOrCreate(&guild, Guild{GuildID: guildID})
	return
}

func (db *Database) GetMemberStats(guildID string, userID string) Member {
	guild, _ := db.GetGuildSettings(guildID)

	return guild.Members[userID]
}

func main() {
	db := &Database{}
	db.Setup()

	guild := &Guild{
		GayConf: GayConfig{
			GayEmote: "Your dad hahahaha",
		},
	}

	db.UpdateGuildSettings("12345", guild)
	tmp, _ := db.GetGuildSettings("12345")
	fmt.Printf("\n\nTHINGY: %v", tmp)
}
