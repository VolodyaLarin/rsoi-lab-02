import express from "express"
import morgan from "morgan"
import bodyParser from "body-parser";

import {TicketService} from "./services/ticketService.js";
import {BonusService} from "./services/bonusService.js";
import {FlightService} from "./services/flightService.js";


const app = express()

app.use(morgan())
app.use(bodyParser())

const router = express.Router()

router.get('/flights', async (req, resp) => {
    const page = req.query.page || "1"
    const size = req.query.size || "100"
    try {
        const flights = await FlightService.GetFlights(page, size);

        resp.json(flights);
        return
    } catch (err) {
        resp.sendStatus(502)
    }
})


const ticketsFill = async (tickets) => {
    const flightsNumbers = [...new Set(tickets.map(x => x.flightNumber))]

    const flights = []

    try {
        const res = await FlightService.GetFlights(1, Math.max(1000, tickets.length), flightsNumbers)
        flights.push(...res.items)
    } catch (err) {
        console.log(err)
    }

    return tickets.map((x) => {
        const flight = flights.find(z => z.flightNumber === x.flightNumber) || {
            fromAirport: "A1", toAirport: "A2", date: '2000-02-02T20:00:00Z'
        }
        return {
            ...x, ...flight
        }
    })
}

router.get('/tickets', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    try {
        const data = await TicketService.GetTickets(username)
        resp.json(await ticketsFill(data))
    } catch (err) {
        console.log(err)
        resp.sendStatus(502)
    }

    return
})


router.get('/tickets/:uid', async (req, resp) => {
    const username = req.header('X-USER-NAME')
    const uid = req.params.uid

    try {
        const data = await TicketService.GetTickets(username, [uid])
        if (data.length === 0) {
            return resp.sendStatus(404)
        }
        resp.json((await ticketsFill([data[0]]))[0])

    } catch (err) {
        console.log(err)
        resp.sendStatus(502)
    }

    return
})


router.get('/me', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    let privilege = {
        balance: 0, status: "BRONZE"
    }
    let tickets = []

    try {
        const data = await BonusService.GetBonusDetails(username)
        privilege.status = data.status
        privilege.balance = data.balance
    } catch (err) {
        console.log(err)
    }
    try {
        const data = await TicketService.GetTickets(username)
        tickets = await ticketsFill(data)
    } catch (err) {
        console.log(err)
    }


    resp.json({
        tickets, privilege
    })
})


router.get('/privilege', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    try {
        const data = await BonusService.GetBonusDetails(username)
        resp.send(data)
    } catch (err) {
        console.log(err)
        resp.send(502)
    }
})


router.post('/tickets', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    const flightNumber = req.body.flightNumber
    const price = req.body.price
    const paidFromBalance = req.body.paidFromBalance

    let ticket = null
    try {
        ticket = await TicketService.CreateTicket(flightNumber, price, username)
    } catch (err) {
        console.log(err)
        resp.send(500)
        return
    }
    let bonusesItem = null;
    try {
        bonusesItem = await BonusService.CreateBonusItem(flightNumber, price, paidFromBalance, ticket.ticketUid, username)
    } catch (err) {
        console.log(err);

        await TicketService.DeleteTicket(ticket.ticketUid, username)

        resp.sendStatus(500)
        return
    }

    console.log(bonusesItem)
    const bonuses = bonusesItem.item.operationType === 'DEBIT_THE_ACCOUNT' ? bonusesItem.item.balanceDiff : 0

    const ticketDetailed = (await ticketsFill([ticket]))[0]

    const paidByMoney = price - bonuses
    const paidByBonuses = bonuses

    resp.status(200).json({
        ...ticketDetailed, paidByBonuses, paidByMoney, privilege: bonusesItem.privelege
    })

    return
})


router.delete('/tickets/:uid', async (req, resp) => {
    const username = req.header('X-USER-NAME')
    const uid = req.params.uid

    try {
        await TicketService.DeleteTicket(uid, username)
    } catch (err) {
        console.log(err)

        resp.sendStatus(500)
        return
    }

    try {
        await BonusService.DeleteBonusItem(uid, username)
    } catch {
        resp.sendStatus(500)
        return
    }
    resp.sendStatus(204)

    return
})

app.use('/api/v1/', router)

app.listen(8080)
