import {
  TelemetryRecord,
  TelemetryResponse,
  AnomalyRecord,
} from "../types/telemetry";

const API_BASE_URL = process.env.REACT_APP_API_URL || "/api/v1";
const WS_URL = process.env.REACT_APP_WS_URL || "ws://localhost:3000/ws";

export class TelemetryService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;

  async getCurrentTelemetry(): Promise<TelemetryRecord> {
    const response = await fetch(`${API_BASE_URL}/telemetry/current`);
    if (!response.ok) {
      throw new Error("Failed to fetch current telemetry");
    }
    return response.json();
  }

  async getTelemetryHistory(
    startTime: string,
    endTime: string,
    page: number = 1,
    limit: number = 20
  ): Promise<TelemetryResponse> {
    const response = await fetch(
      `${API_BASE_URL}/telemetry?start_time=${startTime}&end_time=${endTime}&page=${page}&limit=${limit}`
    );
    if (!response.ok) {
      throw new Error("Failed to fetch telemetry history");
    }
    return response.json();
  }

  async getAnomalies(
    startTime: string,
    endTime: string,
    page: number = 1,
    limit: number = 20
  ): Promise<TelemetryResponse> {
    const response = await fetch(
      `${API_BASE_URL}/telemetry/anomalies?` +
        `start_time=${startTime}&` +
        `end_time=${endTime}&` +
        `page=${page}&` +
        `page_size=${limit}`
    );
    if (!response.ok) {
      throw new Error("Failed to fetch anomalies");
    }
    return response.json();
  }

  connectWebSocket(onMessage: (data: TelemetryRecord) => void): void {
    if (this.ws) {
      this.ws.close();
    }

    this.ws = new WebSocket(WS_URL);

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      onMessage(data);
    };

    this.ws.onclose = () => {
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        setTimeout(() => {
          this.reconnectAttempts++;
          this.connectWebSocket(onMessage);
        }, 1000 * Math.pow(2, this.reconnectAttempts));
      }
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  disconnectWebSocket(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.reconnectAttempts = 0;
  }
}
