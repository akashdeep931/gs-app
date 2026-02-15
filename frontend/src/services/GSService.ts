import type { ShipmentResponse } from "../types/shipment";
import { ValidationError, NetworkError } from "../types/errors";

export class GSService {
  private static instance: GSService;
  private readonly baseUrl: string;

  private constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  static getInstance(
    baseUrl: string = "http://localhost:8000/v1",
  ): GSService {
    if (!GSService.instance) {
      GSService.instance = new GSService(baseUrl);
    }
    return GSService.instance;
  }

  validateQuantity(value: number): string | null {
    if (!Number.isInteger(value)) return "Quantity must be a whole number";
    if (value <= 0) return "Quantity must be greater than 0";
    if (value > 1000000) return "Quantity cannot exceed 1000000";
    return null;
  }

  async calculateShipment(items: number): Promise<ShipmentResponse> {
    const error = this.validateQuantity(items);
    if (error) throw new ValidationError(error);

    let res: Response;

    try {
      res = await fetch(`${this.baseUrl}/shipment`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ items }),
      });
    } catch {
      throw new NetworkError("Unable to reach the server. Please try again.");
    }

    if (!res.ok) {
      const body = await res.text().catch(() => "");
      throw new NetworkError(body || `Request failed (${res.status})`);
    }

    return res.json() as Promise<ShipmentResponse>;
  }
}
