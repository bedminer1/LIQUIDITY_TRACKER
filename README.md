# StableTide

StableTide is a powerful tool designed to monitor and assess liquidity risks for various financial assets, including tokenized assets and ETFs. It leverages AI-powered predictions to identify potential liquidity shortfalls and provide actionable insights for risk management.

## About

StableTide integrates a Go backend, a Python-based AI microservice, and a user-friendly front-end dashboard to streamline the monitoring and prediction of liquidity metrics. It works by analyzing bid-ask spreads, trading volumes, and other market data to assess potential risks and generate predictive insights.

### Key Features:
- **Real-Time(TODO) Data Analysis**: Fetch and process records for various assets.
- **AI-Driven Predictions**: Utilize LSTM models to forecast future liquidity conditions.
- **Customizable Reports**: Generate detailed reports on current and predicted liquidity risks.

## How to Use

### API Endpoints

The StableTide backend provides the following endpoints for interacting with the system:

1. **Fetch Records**:
   - **Path**: `/records`
   - **Description**: Retrieve historical market data for specified assets.
   - **Query Parameters**:
     - `asset`: Asset type (e.g., "ETF", "Crypto").
     - `start`: Start date in `YYYY-MM-DD` format.
     - `end`: End date in `YYYY-MM-DD` format.

   Example Request:
    GET /records?asset=ETF_EMB&start=2023-01-01&end=2023-12-31

2. **Fetch Predictions**:
- **Path**: `/predictions`
- **Description**: Get AI-generated predictions for specified assets based on historical data.
- **Query Parameters**:
  - `asset`: Asset type (e.g., "ETF", "Crypto").
  - `time_to_predict`: Duration to predict into the future (in seconds).
  - `time_interval`: Time interval between predictions (in seconds).

Example Request:
    GET /predictions?asset=ETF_HYG&time_to_predict=3600&time_interval=300


3. **Generate Reports**:
- **Path**: `/report`
- **Description**: Generate a comprehensive report on liquidity risks, combining current data and AI predictions.
- **Query Parameters**:
  - `asset`: Asset type (e.g., "ETF", "Crypto").
  - `start`: Start date in `YYYY-MM-DD` format.
  - `end`: End date in `YYYY-MM-DD` format.

Example Request:
    GET /report?asset=ETF_LQD&start=2023-01-01&end=2023-12-31

### Running the Application

1. **Backend Setup**:
- Ensure Go is installed.
- Run the backend server:
  ```bash
  cd backend/cmd/api
  go run .
  ```

2. **AI Microservice**:
- Ensure Python is installed.
- Start the microservice:
  ```bash
  python ai_microservice.py
  ```

3. **Frontend**:
- Navigate to the `frontend` folder.
- Run the development server:
  ```bash
  npm run dev
  ```

### Requirements

- Python 3.10-3.11
- Go 1.19+
- Node.js 16+
- TensorFlow and supporting libraries installed in the Python environment.

---