# Step 2: Installing Dependencies
With the virtual environment active, you will use Python's package manager to download and install only two external tools (libraries):

The MQTT Client: A library that will teach Python how to speak the Mosquitto protocol so it can subscribe and listen to messages. **DONE**

The PostgreSQL Adapter: A library that will act as the "driver" or translator so Python can open the door to your database and execute SQL statements. 

# Step 3: File Structuring **DONE**
You will create a new blank file in your project (for example, bridge_service.py) and open it in your text editor. At the top of this file, you must import the two libraries you installed in the previous step, along with Python's native tools for handling timestamps and JSON files. 

# Step 4: Credentials Configuration (Global Variables) **DONE**
Next, you will define in plain text all the necessary "keys" to connect to both systems.

For the Database: You will define the server (localhost), the port (5432), the database name, the user, and the password.

For the Broker: You will define the Mosquitto IP (localhost), its port (1883), and, very importantly, the path or topic you are going to subscribe to (the telemetry path using a wildcard to listen to all nodes at the same time).

# Step 5: Database Connection Logic
You will create a function exclusively dedicated to attempting to open the connection with PostgreSQL using the credentials from the previous step. This function must include an error prevention block (to alert you if the database is down) and, if successful, it must keep the connection "live" and available for the rest of the program.

# Step 6: Setting Up Reaction Events (Callbacks)
This is the core part of the program. You must define two functions that you won't call yourself; instead, Mosquitto will execute them automatically when certain things happen:

The Connection Event: An instruction that tells Python: "As soon as you successfully connect to Mosquitto, subscribe immediately to the telemetry path."

The Message Reception Event: The instruction triggered every time an ESP32 sends data. Its workflow must be:

Receive the text packet and translate it from JSON format into a format that Python understands.

Extract the individual variables (node ID, soil moisture, temperature, etc.).

Generate the exact date and time (timestamp) of the server at that very millisecond.

Take the database connection, prepare an INSERT INTO SQL query for your Telemetry table (Hypertable), inject the data, and commit the save to the hard drive.

# Step 7: Main Loop Configuration (The Daemon)
In the final part of your file, you will assemble the "engine" of the service. Here, you must instantiate the MQTT client, assign the two reaction functions you created in step 6, and command it to connect to the broker. Finally, you will execute an "infinite loop" command; this is what turns your script into a continuous service that never stops and keeps listening in the background forever.

# Step 8: Execution and Comprehensive Testing
To test it, you will make sure that your Docker containers (Postgres and Mosquitto) are up and running. Then, you will run your new Python file from the terminal. If everything is correct, the script will remain frozen waiting for data. Using the MQTT Explorer tool, you will send a test JSON simulating a sensor; at that moment, you should see how your Python terminal reacts, processes the data, and silently inserts it into your TimescaleDB database.