package mongo

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//var DocNameTxn = "txn"

type DbSetup struct {
	Host   string
	DbName string
}

type Db struct {
	database *mgo.Database
	session  *mgo.Session
	setup    *DbSetup
}

func (d *Db) Setup(setup *DbSetup) {
	d.setup = setup
}

func (d *Db) Open() error {
	var err error
	if d.session, err = mgo.Dial(d.setup.Host); err != nil {
		return err
	}
	d.database = d.session.DB(d.setup.DbName)
	return nil
}

func (d *Db) C(collectionname string) *mgo.Collection {
	return d.database.C(collectionname)
}

func (d *Db) GridFS(prefix string) *mgo.GridFS {
	return d.database.GridFS(prefix)
}

func (d *Db) Close() {
	d.session.Close()
	d.database = nil
}

func NewDb() *Db {
	var db Db
	return &db
}

func ObjectIdToHexStr(objid bson.ObjectId) string {
	objId := fmt.Sprintf("%s", objid)
	objId = strings.Replace(objId, "ObjectIdHex(\"", "", -1)
	objId = strings.Replace(objId, "\")", "", -1)
	return objId
}

func ObjectIdToHexStrByInterface(objid interface{}) string {
	objId := fmt.Sprintf("%s", objid)
	objId = strings.Replace(objId, "ObjectIdHex(\"", "", -1)
	objId = strings.Replace(objId, "\")", "", -1)
	return objId
}

func ObjectIdHex(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

func NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}
