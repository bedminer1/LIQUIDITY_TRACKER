import sqlite3
import pandas as pd
import numpy as np
from sklearn.preprocessing import MinMaxScaler
from tensorflow.keras.models import Sequential, load_model
from tensorflow.keras.layers import LSTM, Dense, Masking, TimeDistributed # Added TimeDistributed Layer
from sklearn.metrics import mean_squared_error
from sklearn.model_selection import train_test_split
import os
from tensorflow.keras.preprocessing.sequence import pad_sequences

def load_data(db_file):
    conn = sqlite3.connect(db_file)
    query = "SELECT * FROM records"
    data = pd.read_sql_query(query, conn)
    conn.close()
    data['timestamp'] = pd.to_datetime(data['timestamp'], format='%Y-%m-%d %H:%M:%S+00:00')
    return data

def create_sequences(data, future_steps=100):
    sequences = []
    labels = []
    for i in range(future_steps, len(data) - future_steps + 1):
        seq = data[:i]
        label = data[i + future_steps - 1]
        sequences.append(seq)
        labels.append(label)
    return sequences, labels

def pad_and_prepare_sequences(sequences, labels):
    X_padded = pad_sequences(sequences, padding='pre', dtype='float32')
    Y = np.array(labels)
    return X_padded, Y

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
            LSTM(50, return_sequences=False), # Modified this line: the last LSTM layer should not return sequences.
            Dense(3)  # No more TimeDistributed Layer. Output is now a vector
        ])
        model.compile(optimizer='adam', loss='mean_squared_error')
        print("Created a new model.")
    
    model.fit(X_train, Y_train, validation_data=(X_val, Y_val), epochs=epochs, batch_size=batch_size, verbose=1)
    
    predictions = model.predict(X_test)
    mse = mean_squared_error(Y_test, predictions) # Modified this line: no reshape
    print(f"Model Mean Squared Error on Test Data: {mse}")
    
    model.save(model_file)
    print(f"Model saved as '{model_file}'.")

train_model(r"..\backend\market_data.db")
