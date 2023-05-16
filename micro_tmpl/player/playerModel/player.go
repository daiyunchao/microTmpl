package playerModel

type Player struct {
	Id       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}
