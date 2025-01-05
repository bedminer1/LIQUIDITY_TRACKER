from flask import Flask, request, jsonify
import pandas as pd
import tensorflow as tf
from tensorflow.keras.models import load_model
from sklearn.preprocessing import MinMaxScaler
import numpy as np

# Initialize Flask app
app = Flask(__name__)

# Load the trained model
model = load_model("trained_model.keras")

# Endpoint to handle predictions
@app.route("/predict", methods=["POST"])
def predict():
    # Parse incoming JSON data
    data = request.get_json()
    df = pd.DataFrame(data)
    asset_type = data[0]["asset_type"]

    # Ensure required columns are present
    required_columns = ["timestamp", "bid_ask_spread", "volume", "bid_price"]
    if not all(col in df.columns for col in required_columns):
        missing = [col for col in required_columns if col not in df.columns]
        return jsonify({"error": f"Missing required columns: {missing}"}), 400

    # Define the desired window size (can be increased for better accuracy)
    window_size = 20  # Increase from 10 to 20 for more historical context

    # Handle cases where fewer records are sent than the window size
    if len(df) < window_size:
        # Pad with duplicate rows of the first record to meet window size
        padding = [df.iloc[0].to_dict()] * (window_size - len(df))
        df = pd.concat([pd.DataFrame(padding), df], ignore_index=True)

    # Use the last `window_size` rows as the input for predictions
    last_known_data = df[["bid_ask_spread", "volume", "bid_price"]].iloc[-window_size:].values

    # Scale the input data using the same scaler from training
    scaler.fit(last_known_data)  # Fit scaler on the current batch of data
    last_known_data_scaled = scaler.transform(last_known_data)
    future_data = [last_known_data_scaled]


   # Generate future timestamps
    start_time = pd.to_datetime(df["timestamp"].iloc[-1])
    time_interval = data[0].get("time_interval", 86400)  # Default interval in seconds
    time_to_predict = data[0].get("time_to_predict", 604800)  # Default prediction duration (1 week)
    future_timestamps = [
        start_time + pd.to_timedelta(i * time_interval, unit="s")
        for i in range(1, int(time_to_predict / time_interval) + 1)
    ]

    # Generate predictions sequentially
    predictions = []
    for _ in range(len(future_timestamps)):
        pred_scaled = model.predict(np.expand_dims(future_data[-1], axis=0))  # Predict with scaled data
        predictions.append(pred_scaled[0])
        # Append new prediction to future_data and maintain rolling window
        future_data.append(np.vstack([future_data[-1][1:], pred_scaled[0]]))

    # Rescale predictions back to the original scale
    predictions_rescaled = scaler.inverse_transform(predictions)

    predicted_records = []
    for timestamp, pred in zip(future_timestamps, predictions_rescaled):
        predicted_records.append({
            "asset_type": asset_type,
            "timestamp": timestamp.isoformat(),
            "bid_ask_spread": float(pred[0]),
            "volume": float(pred[1]),
            "bid_price": float(pred[2])
        })

    return jsonify(predicted_records)


# Run the microservice
if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5433)