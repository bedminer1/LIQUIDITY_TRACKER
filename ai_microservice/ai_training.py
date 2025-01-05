import sqlite3  # To connect to our .db file
import pandas as pd  # To manage our data easily like a table
import numpy as np  # For numerical operations
from sklearn.preprocessing import MinMaxScaler  # To scale data for LSTM
from tensorflow.keras.models import Sequential, load_model  # To build and load the LSTM model
from tensorflow.keras.layers import LSTM, Dense  # LSTM layers
from sklearn.metrics import mean_squared_error  # To evaluate model accuracy
from sklearn.model_selection import train_test_split  # To split data into training, testing, and validation sets
import os  # To check if a file exists

# ===== Section 1: Train and Save the Model =====

# Function to load data from the database
def load_data(db_file):
    conn = sqlite3.connect(db_file)  # Connect to the database
    query = "SELECT * FROM records"  # SQL to get everything from the table
    data = pd.read_sql_query(query, conn)  # Read the data into a DataFrame
    conn.close()  # Close the connection
    # Convert timestamp to datetime format
    data['timestamp'] = pd.to_datetime(data['timestamp'], format='%Y-%m-%dT%H:%M:%SZ')
    return data

# Function to create sequences for the LSTM model
def create_sequences(data, sequence_length):
    sequences = []
    for i in range(len(data) - sequence_length):
        seq = data[i:i + sequence_length]
        label = data[i + sequence_length]
        sequences.append((seq, label))
    return np.array(sequences, dtype=object)

# Function to train the model
def train_model(db_file, model_file, sequence_length, epochs, batch_size):
    # Step 1: Load the data
    data = load_data(db_file)

    # Step 2: Scale the data
    scaler = MinMaxScaler()
    data_scaled = scaler.fit_transform(data[['bid_ask_spread', 'volume', 'bid_price']])
    data_scaled = pd.DataFrame(data_scaled, columns=['bid_ask_spread', 'volume', 'bid_price'])
    data_scaled['timestamp'] = data['timestamp']
    data_scaled['asset_type'] = data['asset_type']

    # Step 3: Prepare sequences for training
    all_sequences = []
    for asset_type in data_scaled['asset_type'].unique():
        asset_data = data_scaled[data_scaled['asset_type'] == asset_type][['bid_ask_spread', 'volume', 'bid_price']].values
        sequences = create_sequences(asset_data, sequence_length)
        all_sequences.extend(sequences)

    # Step 4: Split sequences into inputs (X) and labels (Y)
    X, Y = zip(*all_sequences)
    X = np.array(X)
    Y = np.array(Y)

    # Step 5: Split the data into training, testing, and validation sets
    X_train, X_temp, Y_train, Y_temp = train_test_split(X, Y, test_size=0.4, random_state=42)
    X_test, X_val, Y_test, Y_val = train_test_split(X_temp, Y_temp, test_size=0.5, random_state=42)

    # Step 6: Check if there's an existing model to continue training
    if os.path.exists("trained_model.keras"):
        model = load_model(model_file)
        print("Loaded existing model.")
    else:
        # Build a new model
        model = Sequential()
        model.add(LSTM(50, return_sequences=True, input_shape=(X_train.shape[1], X_train.shape[2])))
        model.add(LSTM(50))
        model.add(Dense(3))  # Predicting bid_ask_spread, volume, and bid_price
        model.compile(optimizer='adam', loss='mean_squared_error')
        print("Created a new model.")

    # Step 7: Train the model
    model.fit(X_train, Y_train, validation_data=(X_val, Y_val), epochs=epochs, batch_size=batch_size, verbose=1)

    # Step 8: Evaluate the model on the test set
    predictions = model.predict(X_test)
    mse = mean_squared_error(Y_test, predictions)
    print(f"Model Mean Squared Error on Test Data: {mse}")

    # Step 9: Save the model in Keras .h5 format
    model.save('trained_model.keras')
    print("Model saved as 'trained_model.keras'.")

train_model(r"..\backend\market_data.db", model_file='trained_model.keras', sequence_length=10, epochs=10, batch_size=32)