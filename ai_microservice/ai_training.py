import sqlite3  # To connect to our .db file
import pandas as pd  # To manage our data easily like a table
import numpy as np  # For numerical operations
from sklearn.preprocessing import MinMaxScaler  # To scale data for LSTM
from tensorflow.keras.models import Sequential, load_model  # To build and load the LSTM model
from tensorflow.keras.layers import LSTM, Dense, Masking  # LSTM layers and Masking layer
from sklearn.metrics import mean_squared_error  # To evaluate model accuracy
from sklearn.model_selection import train_test_split  # To split data into training, testing, and validation sets
import os  # To check if a file exists
from tensorflow.keras.preprocessing.sequence import pad_sequences  # For padding sequences

# ===== Section 1: Train and Save the Model =====

# Function to load data from the database
def load_data(db_file):
    conn = sqlite3.connect(db_file)  # Connect to the database
    query = "SELECT * FROM records"  # SQL to get everything from the table
    data = pd.read_sql_query(query, conn)  # Read the data into a DataFrame
    conn.close()  # Close the connection
    # Convert timestamp to datetime format
    data['timestamp'] = pd.to_datetime(data['timestamp'], format='%Y-%m-%d %H:%M:%S+00:00')
    return data

# Function to create sequences for the LSTM model
def create_sequences(data):
    sequences = []
    for i in range(1, len(data)):
        seq = data[:i]  # Take all data points from the beginning up to the current point
        label = data[i]  # The next data point to predict
        sequences.append((seq, label))
    return np.array(sequences, dtype=object)

# Function to pad sequences to the same length
def pad_and_prepare_sequences(sequences):
    X, Y = zip(*sequences)
    X_padded = pad_sequences(X, padding='pre', dtype='float32')  # Pad sequences with zeros at the start
    Y = np.array(Y)
    return X_padded, Y

# Function to train the model
def train_model(db_file, model_file, epochs, batch_size):
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
        sequences = create_sequences(asset_data)
        all_sequences.extend(sequences)

    # Step 4: Pad sequences and split into inputs (X) and labels (Y)
    X, Y = pad_and_prepare_sequences(all_sequences)

    # Step 5: Split the data into training, testing, and validation sets
    X_train, X_temp, Y_train, Y_temp = train_test_split(X, Y, test_size=0.4, random_state=42)
    X_test, X_val, Y_test, Y_val = train_test_split(X_temp, Y_temp, test_size=0.5, random_state=42)

    # Step 6: Check if there's an existing model to continue training
    if os.path.exists("trained_model.keras"):
        model = load_model(model_file)
        print("Loaded existing model.")
    else:
        # Build a new model with a Masking layer
        model = Sequential()
        model.add(Masking(mask_value=0.0, input_shape=(None, X_train.shape[2])))  # Masking layer to ignore padded zeros
        model.add(LSTM(50, return_sequences=True))
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

train_model(r"..\backend\market_data.db", model_file='trained_model.keras', epochs=10, batch_size=32)
