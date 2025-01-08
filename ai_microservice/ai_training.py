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

# ===== Section 1: Load Data =====
def load_data(db_file):
    conn = sqlite3.connect(db_file)
    query = "SELECT * FROM records"
    data = pd.read_sql_query(query, conn)
    conn.close()
    data['timestamp'] = pd.to_datetime(data['timestamp'], format='%Y-%m-%d %H:%M:%S+00:00')
    return data

# ===== Section 2: Create Multi-Step Sequences =====
def create_sequences(data, future_steps=100):
    sequences = []
    labels = []
    for i in range(len(data) - future_steps):
        seq = data[:i+1]  # Take all data points up to the current point
        label = data[i+1:i+1+future_steps, :]  # The next `future_steps` points to predict
        sequences.append(seq)
        labels.append(label)
    return sequences, labels

# ===== Section 3: Pad Sequences =====
def pad_and_prepare_sequences(sequences, labels):
    X_padded = pad_sequences(sequences, padding='pre', dtype='float32')
    Y = np.array(labels)
    return X_padded, Y

# ===== Section 4: Train the Model =====
def train_model(db_file, model_file='trained_model.keras', epochs=10, batch_size=32, future_steps=100):
    data = load_data(db_file)
    
    scaler = MinMaxScaler()
    scaled_data = scaler.fit_transform(data[['bid_ask_spread', 'volume', 'bid_price']])
    data_scaled = pd.DataFrame(scaled_data, columns=['bid_ask_spread', 'volume', 'bid_price'])
    
    all_sequences, all_labels = create_sequences(data_scaled.values, future_steps)
    X, Y = pad_and_prepare_sequences(all_sequences, all_labels)
    
    X_train, X_temp, Y_train, Y_temp = train_test_split(X, Y, test_size=0.4, random_state=42)
    X_test, X_val, Y_test, Y_val = train_test_split(X_temp, Y_temp, test_size=0.5, random_state=42)
    
    if os.path.exists(model_file):
        model = load_model(model_file)
        print("Loaded existing model.")
    else:
        model = Sequential([
            Masking(mask_value=0.0, input_shape=(None, X_train.shape[2])),
            LSTM(50, return_sequences=True),
            LSTM(50),
            Dense(3)  # Predicting 100 sets of bid_ask_spread, volume, and bid_price
        ])
        model.compile(optimizer='adam', loss='mean_squared_error')
        print("Created a new model.")
    
    model.fit(X_train, Y_train, validation_data=(X_val, Y_val), epochs=epochs, batch_size=batch_size, verbose=1)
    
    predictions = model.predict(X_test)
    mse = mean_squared_error(Y_test.reshape(-1, 3), predictions.reshape(-1, 3))
    print(f"Model Mean Squared Error on Test Data: {mse}")
    
    model.save(model_file)
    print(f"Model saved as '{model_file}'.")

# ===== Run the Training =====
train_model(r"..\backend\market_data.db")
