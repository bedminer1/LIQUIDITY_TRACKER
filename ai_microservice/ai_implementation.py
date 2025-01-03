import sqlite3  # To connect to our .db file
import pandas as pd  # To manage our data easily like a table
import numpy as np  # For numerical operations
import matplotlib.pyplot as plt  # To make pretty graphs
from sklearn.preprocessing import MinMaxScaler  # To scale data for LSTM
from tensorflow.keras.models import Sequential, load_model  # To build and load the LSTM model
from tensorflow.keras.layers import LSTM, Dense  # LSTM layers
import os  # To check if a file exists
# ===== Section 2: Load and Use the Model to Predict =====

def predict_and_plot(db_file, model_file, time_to_predict, time_interval):
    # Step 1: Load the model
    if not os.path.exists(model_file):
        raise FileNotFoundError(f"Model file '{model_file}' not found.")
    model = load_model(model_file)
    print("Loaded model from file.")

    # Step 2: Load the data
    data = load_data(db_file)

    # Step 3: Scale the data
    scaler = MinMaxScaler()
    data_scaled = scaler.fit_transform(data[['timestamp', 'bid_ask_spread', 'volume', 'bid_price']])

    # Step 4: Prepare the last known data as input for prediction
    last_known_data = data_scaled[-10:]  # Take the last 10 rows as input
    future_data = [last_known_data]

    # Step 5: Generate future timestamps based on the given time interval
    start_time = data['timestamp'].iloc[-1]
    future_timestamps = [start_time + i * time_interval for i in range(1, int(time_to_predict / time_interval) + 1)]

    # Step 6: Predict future values
    predictions = []
    for _ in range(len(future_timestamps)):
        pred = model.predict(np.array(future_data))
        predictions.append(pred[0])
        future_data.append(pred)
        future_data = future_data[1:]

    # Step 7: Rescale predictions back to original values
    predictions = scaler.inverse_transform(predictions)

    # Step 8: Convert predictions to a DataFrame for easier plotting
    predictions_df = pd.DataFrame(predictions, columns=['bid_ask_spread', 'volume', 'bid_price'])
    predictions_df['timestamp'] = future_timestamps

    # Step 9: Plot the predictions
    plt.figure(figsize=(10, 6))
    plt.plot(data['timestamp'], data['bid_ask_spread'], label='Historical Bid-Ask Spread')
    plt.plot(predictions_df['timestamp'], predictions_df['bid_ask_spread'], label='Predicted Bid-Ask Spread')
    plt.xlabel('Timestamp')
    plt.ylabel('Bid-Ask Spread')
    plt.title('Bid-Ask Spread Prediction')
    plt.legend()
    plt.show()

    plt.figure(figsize=(10, 6))
    plt.plot(data['timestamp'], data['volume'], label='Historical Volume')
    plt.plot(predictions_df['timestamp'], predictions_df['volume'], label='Predicted Volume')
    plt.xlabel('Timestamp')
    plt.ylabel('Volume')
    plt.title('Volume Prediction')
    plt.legend()
    plt.show()

    plt.figure(figsize=(10, 6))
    plt.plot(data['timestamp'], data['bid_price'], label='Historical Bid Price')
    plt.plot(predictions_df['timestamp'], predictions_df['bid_price'], label='Predicted Bid Price')
    plt.xlabel('Timestamp')
    plt.ylabel('Bid Price')
    plt.title('Bid Price Prediction')
    plt.legend()
    plt.show()

# Example usage
# train_model("path_to_training_db.db")  # Train a new model or continue training an existing one
# predict_and_plot("path_to_unknown_asset.db", "trained_model.h5", 3600, 300)  # Predict for 1 hour with 5-minute intervals
