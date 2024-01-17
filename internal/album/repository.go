package album

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"go.mongodb.org/mongo-driver/bson/primitive" 
)

// // Define the Platform struct
// type Platform struct {
// 	AmazonMusic bool `json:"amazonMusic"`
// 	JioSaavn    bool `json:"jioSaavn"`
// 	Gaana       bool `json:"gaana"`
// 	// Add other platforms as needed
// }
type Platform struct {
	AmazonMusic map[string]bool `bson:"AmazonMusic"`
	JioSaavn     map[string]bool `bson:"JioSaavn"`
	Gaana        map[string]bool `bson:"Gaana"`
}

type Album struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AlbumID    string             `bson:"album_id" json:"album_id,omitempty"`
    Title       string    `json:"_id,omitempty" bson:"songtitle,omitempty"`
	// Artist      string   `json:"artist"`
	// Producer    string   `json:"producer"`
	// Writer      string   `json:"writer"`
	// ReleaseDate string   `json:"releaseDate"`
	Platforms   Platform `json:"platforms"`
	Producer  string             `bson:"producer"`
	ReleaseDate string           `bson:"releaseDate"`
	Writer    string             `bson:"writer"`
	Artist    string             `bson:"artist"`
}


// AlbumRepository defines the methods for interacting with the data store
type AlbumRepository interface {
	CreateAlbum(ctx context.Context, album Album) error
	GetAlbum(ctx context.Context, albumID string) (Album, error)
	UpdateAlbum(ctx context.Context, albumID string, update Album) error
	DeleteAlbum(ctx context.Context, albumID string) error
	GetSongsWithCapitalTitles(ctx context.Context) ([]string, error)
	SearchAlbums(ctx context.Context, searchTerm string) ([]Album, error)
}

// mongoRepository implements AlbumRepository for MongoDB
type mongoRepository struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

// NewMongoRepository creates a new instance of MongoRepository
func NewMongoRepository(client *mongo.Client, dbName, collectionName string) AlbumRepository {
	return &mongoRepository{client: client, dbName: dbName, collectionName: collectionName}
}

func (r *mongoRepository) CreateAlbum(ctx context.Context, album Album) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	_, err := collection.InsertOne(ctx, album)
	return err
}

func (r *mongoRepository) GetAlbum(ctx context.Context, albumID string) (Album, error) {
	log.Printf("Attempting to retrieve album with ID: %s", albumID)

	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	// filter := bson.M{"_id": albumID}
	filter := bson.M{"album_id": albumID}
	var album Album
	err := collection.FindOne(ctx, filter).Decode(&album)
	if err == mongo.ErrNoDocuments {
		log.Printf("Album with ID %s not found", albumID)

		return Album{}, errors.New("album not found")
	}
	    log.Printf("Successfully retrieved album: %+v", album)

	return album, err
}

func (r *mongoRepository) UpdateAlbum(ctx context.Context, albumID string, update Album) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	filter := bson.M{"_id": albumID}
	updateDoc := bson.M{"$set": update}
	_, err := collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

func (r *mongoRepository) DeleteAlbum(ctx context.Context, albumID string) error {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
	filter := bson.M{"_id": albumID}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}


func (r *mongoRepository) GetSongsWithCapitalTitles(ctx context.Context) ([]string, error) {
	collection := r.client.Database(r.dbName).Collection(r.collectionName)
    // filter := bson.M{"songtitle": bson.M{"$regex": "^[A-Z]+$", "$options": "i"}}

	// filter := bson.M{"songtitle": bson.M{"$regex": "^[A-Z]+$"}}
	filter := bson.M{"songtitle": bson.M{"$regex": "^(?![a-z])[A-Z]"}} 
	//filter := bson.M{"songtitle": bson.M{"$regex": "^[a-z]+$"}} lower


	// Log the filter being used
	log.Printf("Filter used for GetSongsWithCapitalTitles: %v", filter)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		// Log the error
		log.Printf("Error finding documents with capital titles: %s", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	var result []string
	for cur.Next(ctx) {
		var album Album
		if err := cur.Decode(&album); err != nil {
			// Log the decoding error
			log.Printf("Error decoding document: %s", err.Error())
			return nil, err
		}
		result = append(result, album.Title)  // Use album.SongTitle instead of album.Title
	}

	if err := cur.Err(); err != nil {
		// Log the cursor error
		log.Printf("Error iterating over cursor: %s", err.Error())
		return nil, err
	}

	// Log the result
	log.Printf("Songs with capital titles: %v", result)

	return result, nil
}
func (r *mongoRepository) SearchAlbums(ctx context.Context, searchTerm string) ([]Album, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"songtitle": primitive.Regex{Pattern: searchTerm, Options: "i"}},
			{"releasedate": searchTerm},
		},
	}
	collection := r.client.Database(r.dbName).Collection(r.collectionName)

	cur, err := collection.Find(ctx, filter)  
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var albums []Album
	for cur.Next(ctx) {
		var album Album
		if err := cur.Decode(&album); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil

}