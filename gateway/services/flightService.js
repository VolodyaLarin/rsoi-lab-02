import axios from "axios";
import {axiosLogger, generateQuery} from "./utils.js";
import {CreateBreakerContext, DecorateBreaker} from "./cicuitBreakerDecorator.js";

const baseURL = process.env.FLIGHT_SERVICE_URL || "http://localhost:8060/api/v1/flights"
export const FlightService = {
    _service: new axios.Axios({
        baseURL,
    }), GetFlights: async function (page = 1, size = 100, flightsNumbers = []) {
        const request = await this._service.get(`/?${generateQuery({
            page, size, uids: flightsNumbers
        })}`)
        return JSON.parse(request.data);
    }

}

FlightService._service.interceptors.request.use(axiosLogger)


const ctx = CreateBreakerContext((new URL(baseURL).origin + '/manage/health'))

FlightService.GetFlights = DecorateBreaker(ctx, FlightService.GetFlights.bind(FlightService))
