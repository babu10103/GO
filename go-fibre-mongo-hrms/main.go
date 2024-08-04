package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type Employee struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson"_id, omitempty"`
	Name   string             `json:"name"`
	Age    float64            `json:"age"`
	Salary float64            `json:"salary"`
}

type Response struct {
	Id        string     `json:"id,omitempty"`
	Status    string     `json:"status,omitempty"`
	Message   string     `json:"message,omitempty"`
	Employees []Employee `json:"employees,omitempty"`
}

// GenerateObjectId creates a new ObjectId.
func GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

// ObjectIdToHexString converts an ObjectId to a hexadecimal string.
func ObjectIdToHexString(oid primitive.ObjectID) string {
	return oid.Hex()
}

func GenerateObjectIdFromHex(hex string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(hex)
}

func Connect() (MongoInstance, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return MongoInstance{nil, nil}, fmt.Errorf("failed to load .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PWD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbUser == "" || dbPwd == "" {
		return MongoInstance{nil, nil}, fmt.Errorf("invalid dbHost, dbUser or dbPwd")
	}

	connStr := fmt.Sprintf("mongodb://%s:%s@%s:27017", dbUser, dbPwd, dbHost)
	clientOptions := options.Client().ApplyURI(connStr)

	for i := 0; i < 5; i++ {
		log.Println("MongoDB connection attempt:", i+1)
		client, err := mongo.Connect(context.Background(), clientOptions)

		if err == nil {
			err = client.Ping(context.Background(), nil)
			if err == nil {
				log.Println("Connected to MongoDB!")
				return MongoInstance{client, client.Database(dbName)}, nil
			}
		}
		log.Printf("Attempt %d: MongoDB connection failed with error: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	return MongoInstance{nil, nil}, fmt.Errorf("could not connect to MongoDB")

}

func main() {
	mg, err := Connect()

	if err != nil {
		log.Panic(err)
	}

	defer mg.Client.Disconnect(context.Background())

	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {
		var employees []Employee
		cursor, err := mg.DB.Collection("employees").Find(context.Background(), bson.D{})
		if err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Failed to fetch employees: " + err.Error(),
			})
		}
		if err = cursor.All(context.Background(), &employees); err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Failed to fetch employees: " + err.Error(),
			})
		}
		return c.JSON(Response{
			Employees: employees,
		})
	})

	app.Get("/employee/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		collection := mg.DB.Collection("employees")

		// Convert the ID from the URL parameter
		oid, err := GenerateObjectIdFromHex(id)
		if err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Invalid employee ID: " + err.Error(),
			})
		}

		// Check if the employee exists
		var employee Employee
		err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&employee)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(Response{
					Status:  "Failed",
					Message: "Employee for ID " + id + " not found",
				})
			}
			// Log the error for debugging
			log.Printf("Error finding employee: %v", err)
			return c.Status(500).JSON(Response{
				Status:  "Failed",
				Message: "Error finding employee: " + err.Error(),
			})
		}

		// Return the employee data
		return c.JSON(employee)
	})

	app.Post("/employee", func(c *fiber.Ctx) error {
		var employee Employee

		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Invalid request payload: " + err.Error(),
			})
		}

		employee.ID = GenerateObjectId()
		_, err := mg.DB.Collection("employees").InsertOne(context.Background(), employee)

		if err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Failed to create employee: " + err.Error(),
			})
		}
		return c.JSON(Response{
			Id:      employee.ID.Hex(),
			Status:  "Successful",
			Message: "Employee created successfully",
		})

	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var employee Employee

		// Parse the request body into the employee struct
		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Invalid request payload: " + err.Error(),
			})
		}

		collection := mg.DB.Collection("employees")

		// Convert the ID from the URL parameter
		oid, err := GenerateObjectIdFromHex(id)
		if err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Invalid employee ID: " + err.Error(),
			})
		}

		// Check if the employee exists
		var curEmployee Employee
		err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&curEmployee)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(Response{
					Status:  "Failed",
					Message: "Employee for ID " + id + " not found",
				})
			}
			return c.Status(500).JSON(Response{
				Status:  "Failed",
				Message: "Error finding employee: " + err.Error(),
			})
		}

		// Prepare the update document
		updateFields := bson.M{}

		if employee.Name != "" && employee.Name != curEmployee.Name {
			updateFields["name"] = employee.Name
		}

		if employee.Age != 0 && (curEmployee.Age == 0 || employee.Age != curEmployee.Age) {
			updateFields["age"] = employee.Age
		}

		if employee.Salary != 0 && (curEmployee.Salary == 0 || employee.Salary != curEmployee.Salary) {
			updateFields["salary"] = employee.Salary
		}

		if len(updateFields) == 0 {
			return c.Status(200).JSON(Response{
				Status:  "Successful",
				Message: "No fields to update",
			})
		}

		// Perform the update operation
		update := bson.M{"$set": updateFields}
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
		if err != nil {
			return c.Status(500).JSON(Response{
				Status:  "Failed",
				Message: "Failed to update employee: " + err.Error(),
			})
		}

		// Return success response
		return c.Status(200).JSON(Response{
			Id:      id,
			Status:  "Successful",
			Message: "Employee updated successfully",
		})
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		oid, err := GenerateObjectIdFromHex(id)
		if err != nil {
			return c.Status(400).JSON(Response{
				Status:  "Failed",
				Message: "Invalid employee ID: " + err.Error(),
			})
		}
		_, err = mg.DB.Collection("employees").DeleteOne(context.Background(), bson.M{"_id": oid})
		if err != nil {
			return c.Status(500).JSON(Response{
				Status:  "Failed",
				Message: "Failed to delete employee: " + err.Error(),
			})
		}
		return c.Status(200).JSON(Response{
			Status:  "Successful",
			Message: "Employee deleted successfully",
		})
	})

	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		appPort = "3000"
		log.Printf("Defaulting to port %s", appPort)
	}
	log.Printf("Listening on port %s", appPort)
	log.Fatal(app.Listen(":" + appPort))
}
