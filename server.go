package main

import (
	"bufio"
	"crypto/rand"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type NameDb struct {
	*mgo.Session
}

//https://www.w3.org/International/questions/qa-personal-names

var db = Connect()

func Api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	nm := db.RandomName(true, nil, "us")
	Resonse, _ := template.New("listapi").Parse("{{.}}")
	_ = Resonse.Execute(w, nm)
}

func main() {

	defer db.Close()
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/").Name("Api").HandlerFunc(Api)
	log.Fatal(http.ListenAndServe(":8090", router))

	// files := []string{
	// 	"/go_src/raw/nbs/CA.TXT",
	// 	"/go_src/raw/nbs/CT.TXT",
	// 	"/go_src/raw/nbs/DE.TXT",
	// 	"/go_src/raw/nbs/GA.TXT",
	// 	"/go_src/raw/nbs/IA.TXT",
	// 	"/go_src/raw/nbs/IL.TXT",
	// 	"/go_src/raw/nbs/KS.TXT",
	// 	"/go_src/raw/nbs/LA.TXT",
	// 	"/go_src/raw/nbs/MD.TXT",
	// 	"/go_src/raw/nbs/MI.TXT",
	// 	"/go_src/raw/nbs/MO.TXT",
	// 	"/go_src/raw/nbs/MT.TXT",
	// 	"/go_src/raw/nbs/ND.TXT",
	// 	"/go_src/raw/nbs/NH.TXT",
	// 	"/go_src/raw/nbs/NM.TXT",
	// 	"/go_src/raw/nbs/NY.TXT",
	// 	"/go_src/raw/nbs/OK.TXT",
	// 	"/go_src/raw/nbs/PA.TXT",
	// 	"/go_src/raw/nbs/SC.TXT",
	// 	"/go_src/raw/nbs/TX.TXT",
	// 	"/go_src/raw/nbs/VA.TXT",
	// 	"/go_src/raw/nbs/WA.TXT",
	// 	"/go_src/raw/nbs/WV.TXT",
	// 	"/go_src/raw/nbs/AL.TXT",
	// 	"/go_src/raw/nbs/AZ.TXT",
	// 	"/go_src/raw/nbs/CO.TXT",
	// 	"/go_src/raw/nbs/DC.TXT",
	// 	"/go_src/raw/nbs/FL.TXT",
	// 	"/go_src/raw/nbs/HI.TXT",
	// 	"/go_src/raw/nbs/ID.TXT",
	// 	"/go_src/raw/nbs/IN.TXT",
	// 	"/go_src/raw/nbs/KY.TXT",
	// 	"/go_src/raw/nbs/MA.TXT",
	// 	"/go_src/raw/nbs/ME.TXT",
	// 	"/go_src/raw/nbs/MN.TXT",
	// 	"/go_src/raw/nbs/MS.TXT",
	// 	"/go_src/raw/nbs/NC.TXT",
	// 	"/go_src/raw/nbs/NE.TXT",
	// 	"/go_src/raw/nbs/NJ.TXT",
	// 	"/go_src/raw/nbs/NV.TXT",
	// 	"/go_src/raw/nbs/OH.TXT",
	// 	"/go_src/raw/nbs/OR.TXT",
	// 	"/go_src/raw/nbs/RI.TXT",
	// 	"/go_src/raw/nbs/SD.TXT",
	// 	"/go_src/raw/nbs/TN.TXT",
	// 	"/go_src/raw/nbs/UT.TXT",
	// 	"/go_src/raw/nbs/VT.TXT",
	// 	"/go_src/raw/nbs/WI.TXT",
	// 	"/go_src/raw/nbs/WY.TXT",
	// }
	// for _, f := range files {
	// 	db.NamesByStateParser(f)
	// }

	// db.InsertIfNo()
	// db.Ups(true, true, "mike", "us")
}

func Connect() *NameDb {
	// BUG(r): Rewrite DB code with concurency in mind
	session, err := mgo.Dial("mongodb://30oktphkig:8bohxNFy8u@mongo/admin")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return &NameDb{session}
}

func (self *NameDb) InsertIfNo() {
	// MongoResponseObject := MongoResponse{}
	// pipe := []bson.M{{"$match": bson.M{"_id": bson.ObjectIdHex(RecordingID)}}, {"$sort": bson.M{"_id": 1}}}
	// er := session.DB("plankton").Run(bson.D{
	// 	{"aggregate", "x11gui"},
	// 	{"pipeline", pipe},
	// 	{"allowDiskUse", true},
	// }, &MongoResponseObject)

	// if er != nil {
	// 	panic(er)
	// }
	type Rcr struct {
		Male    bool
		Female  bool
		Country string
		Name    string
	}

	result := bson.M{}
	er := self.Session.DB("plankton").Run(bson.D{
		{"insert", "namecollection"},
		{"ordered", true},
		{"documents", []*Rcr{
			{Male: false, Female: true, Name: "helen", Country: "us"},
			{Male: false, Female: true, Name: "helen", Country: "us"}}},
	}, &result)

	if er != nil {
		panic(er)
	}

}

func (self *NameDb) Ups(male, female bool, name, country string) {
	result := bson.M{}
	q := bson.M{"name": name}
	gender := bson.M{}
	if male {
		gender["male"] = male
	}
	if female {
		gender["female"] = female
	}
	u := bson.M{"$set": gender, "$addToSet": bson.M{"country": country}}
	er := self.Session.DB("plankton").Run(bson.D{
		{"update", "namecollection"},
		{"ordered", true},
		{"updates", []bson.M{
			{
				"q":      q,
				"u":      u,
				"upsert": true,
				"multi":  false}}},
	}, &result)

	if er != nil {
		panic(er)
	}

}

func (self *NameDb) NamesByStateParser(path string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := strings.Split(scanner.Text(), ",")
		male := l[1] == "M"
		female := l[1] == "F"
		name := strings.ToLower(l[3])
		country := "us"
		self.Ups(male, female, name, country)
	}

}

func (self *NameDb) RandomName(male, female interface{}, country string) string {
	result := bson.M{}
	gender := bson.M{}
	if male != nil {
		gender["male"] = male
	} else {
		gender["male"] = bson.M{"$exists": false}
	}
	if female != nil {
		gender["female"] = female
	} else {
		gender["female"] = bson.M{"$exists": false}
	}
	gender["country"] = country
	// gender := bson.M{"male": male, "female": female, "country": country}
	pipe := []bson.M{
		{"$match": gender},
		{"$sort": bson.M{"name": 1}},
		{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}
	er := self.Session.DB("plankton").Run(bson.D{
		{"aggregate", "namecollection"},
		{"pipeline", pipe},
		{"allowDiskUse", true},
		// {"ordered", true},
	}, &result)
	if er != nil {
		panic(er)
	}
	numberofrecords := result["result"].([]interface{})[0].(bson.M)["count"].(int)
	r, e := rand.Int(rand.Reader, big.NewInt(int64(numberofrecords)))
	if e != nil {
		panic(e)
	}
	// db.myCollection.find(query).limit(1).skip(r);
	pipe = []bson.M{
		{"$match": gender},
		{"$sort": bson.M{"name": 1}},
		{"$skip": r.Int64()},
		{"$limit": 1}}
	er = self.Session.DB("plankton").Run(bson.D{
		{"aggregate", "namecollection"},
		{"pipeline", pipe},
		{"allowDiskUse", true},
		// {"ordered", true},
	}, &result)
	if er != nil {
		panic(er)
	}
	resname := result["result"].([]interface{})[0].(bson.M)["name"].(string)
	// spew.Dump(result)
	return strings.Title(resname)
}
