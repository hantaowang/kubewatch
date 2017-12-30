from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import time
import os
import cgi

hostName = "localhost"
hostPort = 9000

events = []

class Event:
 
   def __init__(self, message):
       self.message = message
       self.time = time.asctime()

   def __repr__(self):
       return "[{0}] {1}".format(self.time, self.message)


class MyServer(BaseHTTPRequestHandler):

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
	for e in events:
		self.wfile.write(e)
		self.wfile.write("<br>")

    def do_POST(self):
        # Parse the form data posted
        form = cgi.FieldStorage(
            fp=self.rfile,
            headers=self.headers,
            environ={'REQUEST_METHOD': 'POST',
                     'CONTENT_TYPE': self.headers['Content-Type'],
                     })

        # Begin the response
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()

        events.append(Event(form["event"].value))
	self.wfile.write("done")

myServer = HTTPServer((hostName, hostPort), MyServer)
print("{0} Server STARTS @ {1}:{2}".format(time.asctime(), hostName, hostPort))
print("------------------------------\n")
try:
    myServer.serve_forever()
except KeyboardInterrupt:
    pass

myServer.server_close()
print("\n------------------------------")
print("{0} Server STOPS @ {1}:{2}".format(time.asctime(), hostName, hostPort))

