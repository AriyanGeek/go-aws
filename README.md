# Implement a simple Restful API on AWS Lambda using Golang
<h2>The API accept the following JSON requests : </h2>
<b>GET</b>
For geting a device you should send get request to the following URL with id of device you are looking for:
https://mv7sbb4s5b.execute-api.us-east-1.amazonaws.com/devices/:id
<h6>Examlpe:</h6>
GET https://mv7sbb4s5b.execute-api.us-east-1.amazonaws.com/devices/id1
<br/>
<b>POST</b>
For creating a device you should send post request to the following URL with body that contains (id,name,note,serial,deviceModel):
https://mv7sbb4s5b.execute-api.us-east-1.amazonaws.com/devices
Body: 
{
    "id": "id3",
    "deviceModel": "model3",
    "name": "Sensor3",
    "note": "Testing a sensor.",
    "serial": "C020000102"
}
<br/>
<h4>Or you can test API with Postman collection witch you can find in rep</h4>
<h2>Test coverage: </h2>
My test coverage of this code is 91% that you can find in <a>cover.html</a> file 
