# OK Fuel Prices API Documentation

## 1. Overview
The OK Fuel Prices API provides public access to real-time fuel pricing data for OK gas stations in Denmark.

- **Base URL:** `https://mobility-prices.ok.dk/api/v1/fuel-prices`
- **Method:** `GET`
- **Authentication:** None (Public access)
- **Rate Limiting:** Implemented to ensure stability
- **Caching:** Responses may be cached for up to 2 minutes

## 2. Data Structure
The API returns a JSON object containing an array of gas stations.

### Response Fields (Root)
| Field | Type | Description | Mandatory |
| :--- | :--- | :--- | :--- |
| `items` | Array | List of gas stations | Yes |

### Station Object (`items[]`)
| Field | Type | Description | Mandatory |
| :--- | :--- | :--- | :--- |
| `house_number` | String | House number | Yes |
| `postal_code` | Number | Postal code | Yes |
| `city` | String | City | Yes |
| `coordinates` | Object | Geographic coordinates (lat/lng) | Yes |
| `last_updated_time` | String | Timestamp of latest price update | Yes |
| `prices` | Array | List of fuel prices | Yes |

### Price Object (`prices[]`)
| Field | Type | Description | Mandatory |
| :--- | :--- | :--- | :--- |
| `product_name` | String | Product name (e.g., Blyfri 95) | Yes |
| `price` | Number | Price in DKK (two decimals) | Yes |

## 3. Example Request
```bash
curl -X GET "[https://mobility-prices.ok.dk/api/v1/fuel-prices](https://mobility-prices.ok.dk/api/v1/fuel-prices)"
```
## 4. Example Response
```json
{
  "items": [
    {
      "facility_number": 1234,
      "street": "Eksempelvej",
      "house_number": "10",
      "postal_code": 8000,
      "city": "Aarhus C",
      "coordinates": {
        "latitude": 56.1741,
        "longitude": 9.5515
      },
      "last_updated_time": "2025-10-24T08:30:00Z",
      "prices": [
        {
          "product_name": "Blyfri 95",
          "price": 13.79
        },
        {
          "product_name": "Svovlfri Diesel",
          "price": 12.89
        }
      ]
    }
  ]
}
```
