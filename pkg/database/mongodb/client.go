package mongodb

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"

    "github.com/spiri2/Quests/pkg/constants"
    "github.com/spiri2/Quests/pkg/convert"
    "github.com/spiri2/Quests/pkg/database"
    "github.com/spiri2/Quests/pkg/models"
)

const DefaultTimeout = 20 * time.Second
const InvasionCollection = "invasions"
const QuestCollection = "quests"
const ExpiredQuestCollection = "expiredquests"
const RaidCollection = "raids"

type Client struct {
    db       *mongo.Client
    database string
}

var _ database.Client = (*Client)(nil)

func contextWithTimeout() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), DefaultTimeout)

}

func New(address string, username string, password string, database string) (*Client, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    connectionString := "mongodb://"

    if username != "" && password != "" {
        connectionString += username + ":" + password + "@"
    }

    connectionString += address

    client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
    if err != nil {
        return nil, err
    }

    err = client.Connect(ctx)
    if err != nil {
        return nil, err
    }

    err = client.Ping(ctx, readpref.Primary())

    return &Client{
        db:       client,
        database: database,
    }, err
}

func (client *Client) Init() error {
    return nil
}

func (client *Client) invasions() *mongo.Collection {
    return client.db.Database(client.database).Collection(InvasionCollection)
}

func (client *Client) quests() *mongo.Collection {
    return client.db.Database(client.database).Collection(QuestCollection)
}

func (client *Client) expiredQuests() *mongo.Collection {
    return client.db.Database(client.database).Collection(ExpiredQuestCollection)
}

func (client *Client) raids() *mongo.Collection {
    return client.db.Database(client.database).Collection(RaidCollection)
}

func (client *Client) AddInvasions(invasions ...models.Invasion) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    ordered := true
    _, err := client.invasions().InsertMany(ctx, convert.InvasionSliceToInterfaceSlice(invasions), &options.InsertManyOptions{
        Ordered: &ordered,
    })

    return err
}

func (client *Client) GetAllInvasions() ([]models.Invasion, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    cursor, err := client.invasions().Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    var invasions []models.Invasion
    if err = cursor.All(ctx, &invasions); err != nil {
        return nil, err
    }

    return invasions, nil
}

func (client *Client) DeleteInvasionsByID(ids ...string) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    _, err := client.invasions().DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

    return err
}

func (client *Client) AddExpiredQuests(quests ...models.Quest) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    ordered := true
    _, err := client.expiredQuests().InsertMany(ctx, convert.QuestSliceToInterfaceSlice(quests), &options.InsertManyOptions{
        Ordered: &ordered,
    })

    return err
}

func (client *Client) GetAllExpiredQuests() ([]models.Quest, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    cursor, err := client.expiredQuests().Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    var quests []models.Quest
    if err = cursor.All(ctx, &quests); err != nil {
        return nil, err
    }

    return quests, nil
}

func (client *Client) GetExpiredQuestsByType(questType constants.ServiceType) ([]models.Quest, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    cursor, err := client.expiredQuests().Find(ctx, bson.M{"type": questType})
    if err != nil {
        return nil, err
    }

    var quests []models.Quest
    if err = cursor.All(ctx, &quests); err != nil {
        return nil, err
    }

    return quests, nil
}

func (client *Client) DeleteExpiredQuestsByID(ids ...string) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    _, err := client.expiredQuests().DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

    return err
}

func (client *Client) AddQuests(quests ...models.Quest) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    ordered := true
    _, err := client.quests().InsertMany(ctx, convert.QuestSliceToInterfaceSlice(quests), &options.InsertManyOptions{
        Ordered: &ordered,
    })

    return err
}

func (client *Client) GetAllQuests() ([]models.Quest, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    cursor, err := client.quests().Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    var quests []models.Quest
    if err = cursor.All(ctx, &quests); err != nil {
        return nil, err
    }

    return quests, nil
}

func (client *Client) GetQuestsByType(questType constants.ServiceType) ([]models.Quest, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    cursor, err := client.quests().Find(ctx, bson.M{"type": questType})
    if err != nil {
        return nil, err
    }

    var quests []models.Quest
    if err = cursor.All(ctx, &quests); err != nil {
        return nil, err
    }

    return quests, nil
}

func (client *Client) DeleteQuestsByID(ids ...string) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    _, err := client.quests().DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

    return err
}

func (client *Client) AddRaids(raids ...models.Raid) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    ordered := true
    _, err := client.raids().InsertMany(ctx, convert.RaidSliceToInterfaceSlice(raids), &options.InsertManyOptions{
        Ordered: &ordered,
    })

    return err
}

func (client *Client) GetAllRaids() ([]models.Raid, error) {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    cursor, err := client.raids().Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    var raids []models.Raid
    if err = cursor.All(ctx, &raids); err != nil {
        return nil, err
    }

    return raids, nil
}

func (client *Client) DeleteRaidsByID(ids ...string) error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    _, err := client.raids().DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

    return err
}

func (client *Client) Stop() error {
    ctx, cancel := contextWithTimeout()
    defer cancel()

    return client.db.Disconnect(ctx)
}
