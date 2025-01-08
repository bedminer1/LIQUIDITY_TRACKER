from flask import Flask, request, jsonify
import pandas as pd
import numpy as np
from sklearn.preprocessing import MinMaxScaler
from tensorflow.keras.models import load_model

# Initialize Flask app
app = Flask(__name__)

# Load the trained model
model = load_model("trained_model.keras")

# Initialize a scaler with the same feature ranges used during training
scaler = MinMaxScaler(feature_range=(0, 1))

# Endpoint to handle predictions
@app.route("/predict", methods=["POST"])
def predict():
    # Parse query params
    interval_length = request.args.get("time_interval_length", default=86400, type=int)
    intervals = request.args.get("time_intervals", default=30, type=int)

    # Parse incoming JSON data
    data = request.get_json()
    df = pd.DataFrame(data)
    asset_type = ""
    if len(data) != 0:
        asset_type = data[0]["asset_type"]

    # Ensure required columns are present
    required_columns = ["timestamp", "bid_ask_spread", "volume", "bid_price"]
    if not all(col in df.columns for col in required_columns):
        missing = [col for col in required_columns if col not in df.columns]
        return jsonify({"error": f"Missing required columns: {missing}"}), 400

    # Define the desired window size
    window_size = 20  # Ensure consistent input size for the LSTM

    # Handle cases where fewer records are sent than the window size
    if len(df) < window_size:
        # Pad with duplicate rows of the first record to meet window size
        padding = [df.iloc[0].to_dict()] * (window_size - len(df))
        df = pd.concat([pd.DataFrame(padding), df], ignore_index=True)

    # Extract the last `window_size` rows for predictions
    last_known_data = df[["bid_ask_spread", "volume", "bid_price"]].iloc[-window_size:].values

    # Scale the input data using the same scaler from training
    scaler.fit(last_known_data)  # Fit scaler on the current batch of data
    last_known_data_scaled = scaler.transform(last_known_data)
    future_data = [last_known_data_scaled]

    # Generate future timestamps
    start_time = pd.to_datetime(df["timestamp"].iloc[-1])
    future_timestamps = [
        start_time + pd.to_timedelta(i * interval_length, unit="s")
        for i in range(1, intervals + 1)
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

    # Format the predictions as a list of dictionaries
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
