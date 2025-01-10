# StableTide: Liquidity Risk Management for Tokenized Assets

StableTide is a blockchain-focused liquidity tracker designed for financial institutions to monitor, predict, and mitigate liquidity risks associated with tokenized assets. By leveraging advanced machine learning and blockchain integration, StableTide provides actionable insights, empowering institutions to manage tokenized assets proactively.

---

## Problem Statement and Motivation

### Problem Statement  
Tokenized assets, while revolutionary, pose unique challenges:  
- **Fragmented Data**: Blockchain ecosystems are decentralized, making comprehensive liquidity tracking difficult.  
- **High Volatility**: Tokenized assets often exhibit unpredictable trading patterns, increasing liquidity risks.  
- **Limited Tools**: Current solutions are insufficient for managing liquidity risks in decentralized financial systems.  

### Motivation  
Tokenization is transforming traditional finance by:  
- Expanding accessibility to markets.  
- Increasing liquidity for previously illiquid assets.  
- Leveraging blockchain transparency for trust and accountability.  

However, without robust tools like StableTide, institutions risk systemic failures, regulatory non-compliance, and asset devaluation.

---

## Features

### Liquidity Risk Monitoring  
- Tracks trading volume, bid-ask spread, and transaction frequency.  

### Predictive Analytics  
- Uses statistical modeling and LSTM-based forecasting to predict liquidity shortfalls.  

### Recommendations via OpenAI  
- Generates actionable recommendations tailored to specific tokenized assets using OpenAI's GPT models.

---

## Technical Overview

### Backend  
- **Language**: Go (for performance and scalability).  
- **Endpoints**:  
  - `/recommendations`: Analyzes asset liquidity and provides actionable insights.  

#### `/recommendations` Endpoint  
- **Input**:  
  - Query parameters: `start`, `end`, `asset`, `time_intervals`, and `time_interval_length`.  
- **Process**:  
  - Fetches historical asset data from a database.  
  - Generates liquidity predictions using statistical models.  
  - Assesses risks and sends the data to OpenAI for recommendation generation.  
- **Output**:  
  - Analysis (AI-generated insights).  
  - Historical data and predictions.  
  - Comprehensive liquidity report.  

### Frontend  
- Built with **SvelteKit** for an intuitive user interface.  
- Features interactive graphs for bid-ask spread percentage and trading volume trends.  
- Provides real-time insights into liquidity trends and warnings.

---

## How It Works

1. **Input Data**:  
   Users specify asset type, time frame, and prediction intervals.  

2. **Data Aggregation**:  
   StableTide fetches and organizes blockchain transaction data into a database.  

3. **Prediction Generation**:  
   Statistical or ML-based models predict future liquidity trends.  

4. **Recommendations**:  
   AI analyzes predictions and provides actionable advice tailored to the asset.  

5. **Visualization**:  
   Users access detailed graphs and analyses through the frontend.

---

## Installation and Usage

### Prerequisites  
- **Backend**: Go installed (v1.17+).  
- **Frontend**: Node.js and npm installed.  
- **AI Microservice**: Python 3.11 with required libraries.

### Setting Up  
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/StableTide.git
   cd StableTide
   ```
2. **Initialize the Program**:
   ```bash
   make init_all
   ```
3. **Run the Programs**:
   ```bash
      make run_all
   ```

### Requirements
- Python3.11
- Go 1.18+
- Node 16.x+
- Make
- SQLite
