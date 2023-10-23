import axios from "axios";
import {axiosLogger, generateQuery, genUsernameHeaders} from "./utils.js";
import {CreateBreakerContext, DecorateBreaker} from "./cicuitBreakerDecorator.js";


const baseURL = process.env.TICKET_SERVICE_URL || "http://localhost:8070/api/v1/tickets"
export const TicketService = {
    _service: new axios.Axios({
        baseURL,
    }),

    GetTickets: async function (username, uids = []) {
        const request = await this._service.get(`/?${generateQuery({
            uids: uids
        })}`, genUsernameHeaders(username));

        if (request.status !== 200) {
            throw "Request error"
        }

        return JSON.parse(request.data)
    }, CreateTicket: async function (flightNumber, price, username) {
        const request = await this._service.post('/', JSON.stringify({
            flightNumber, price,
        }), genUsernameHeaders(username));
        return JSON.parse(request.data)
    }, DeleteTicket: async function (uid, username) {
        await this._service.delete(`/${uid}`, genUsernameHeaders(username));
        return
    }
}
TicketService._service.interceptors.request.use(axiosLogger)

const ctx = CreateBreakerContext((new URL(baseURL).origin + '/manage/health'))

TicketService.GetTickets = DecorateBreaker(ctx, TicketService.GetTickets.bind(TicketService))
TicketService.CreateTicket = DecorateBreaker(ctx, TicketService.CreateTicket.bind(TicketService))
TicketService.DeleteTicket = DecorateBreaker(ctx, TicketService.DeleteTicket.bind(TicketService))
