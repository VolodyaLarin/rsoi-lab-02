import express from "express"
import axios from "axios"
import morgan from "morgan"
import bodyParser from "body-parser";

import os from 'os'

const ticketService = new axios.Axios({
    baseURL: process.env.TICKET_SERVICE_URL || "http://localhost:8070/api/v1/tickets",
})

const bonusService = new axios.Axios({
    baseURL: process.env.BONUS_SERVICE_URL || "http://localhost:8050/api/v1/bonus",
})

const flightService = new axios.Axios({
    baseURL: process.env.FLIGHT_SERVICE_URL || "http://localhost:8060/api/v1/flights",
})

axios.interceptors.request.use(request => {
    console.log('Request: ', JSON.stringify(request, null, 2))
    return request
})


const generateQuery = (data) => {
    return Object.entries(data).map((kv) => Array.isArray(kv[0]) ?
        kv[1].map(v => `${encodeURIComponent(kv[0])}[]=${encodeURIComponent(v)}`).join('&')
        : `${encodeURIComponent(kv[0])}=${encodeURIComponent(kv[1])}`).join('&')
}

const genUsernameHeaders = (username) => ({
    headers: {
        'X-USER-NAME': String(username),
        'Content-type': 'application/json'
    }
})


const app = express()

app.use(morgan())
app.use(bodyParser())

const router = express.Router()

router.get('/flights', async (req, resp) => {
    const page = req.query.page || "0"
    const size = req.query.size || "100"
    try {
        const request = await flightService.get(`/?${generateQuery({
            page, size
        })}`)

        resp.status(request.status).json(JSON.parse(request.data));
        return
    } catch (err) {
        console.log(err)
        resp.sendStatus(502)
    }
})


const ticketsFill = async (tickets) => {
    const flightsNumbers = [...new Set(tickets.map(x => x.flightNumber))]

    let flights = []

    try {
        const response = await flightService.get(`/?${generateQuery({
            uids: flightsNumbers
        })}`)
        if (response.status === 200) {
            flights = JSON.parse(response.data).items
        }
    } catch {

    }

    return tickets.map((x) => {
        const flight = flights.find(z => z.flightNumber === x.flightNumber) || {
            fromAirport: "A1",
            toAirport: "A2",
            date: '2000-02-02T20:00:00Z'
        }
        return {
            ...x,
            ...flight
        }
    })
}

router.get('/tickets', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    try {
        const request = await ticketService.get(`/`, genUsernameHeaders(username))
        if (request.status !== 200) {
            return resp.sendStatus(request.status)
        }

        const data = JSON.parse(request.data)


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
        const request = await ticketService.get(`/?${generateQuery({
            uids: [uid]
        })}`, genUsernameHeaders(username))
        if (request.status !== 200) {
            return resp.sendStatus(request.status)
        }
        const data = JSON.parse(request.data)

        if (data.length === 0) {
            return resp.sendStatus(404)
        }

        resp.json((await ticketsFill(data))[0])

    } catch (err) {
        console.log(err)
        resp.sendStatus(502)
    }

    return
})


router.get('/me', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    let privilege = {
        balance: 0,
        status: "BRONZE"
    }
    let tickets = []

    try {
        const request = await bonusService.get(`/`, genUsernameHeaders(username))
        const data = JSON.parse(request.data)
        privilege.status = data.status
        privilege.balance = data.balance
    } catch (err) {
        console.log(err)
    }
    try {
        const request = await ticketService.get(`/`, genUsernameHeaders(username))
        tickets = await ticketsFill(JSON.parse(request.data))
    } catch (err) {
        console.log(err)
    }


    resp.json({
        tickets,
        privilege
    })
})


router.get('/privilege', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    try {
        const request = await bonusService.get(`/`, genUsernameHeaders(username))
        const data = JSON.parse(request.data)

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


    const ticketResp = await ticketService.post('/', JSON.stringify({
        flightNumber,
        price,
    }), genUsernameHeaders(username))

    const ticket = JSON.parse(ticketResp.data)


    const bonusResp = await bonusService.post('/', JSON.stringify({
        flightNumber,
        price,
        paidFromBalance,
        ticketUid: ticket.ticketUid
    }), genUsernameHeaders(username))


    const bonusesItem = JSON.parse(bonusResp.data)

    console.log(bonusesItem)

    const bonuses = bonusesItem.item.operationType === 'DEBIT_THE_ACCOUNT' ? bonusesItem.item.balanceDiff : 0


    const ticketDetailed = (await ticketsFill([ticket]))[0]

    const paidByMoney = price - bonuses
    const paidByBonuses = bonuses

    resp.status(200).json({
        ...ticketDetailed,
        paidByBonuses,
        paidByMoney,
        privilege: bonusesItem.privelege
    })

    return
})

router.delete('/tickets/:uid', async (req, resp) => {
    const username = req.header('X-USER-NAME')

    const uid = req.params.uid

    const ticketResp = await ticketService.delete(`/${uid}`, genUsernameHeaders(username))
    const bonusResp = await bonusService.delete(`/${uid}`, genUsernameHeaders(username))

    resp.sendStatus(204)

    return
})

app.use('/api/v1/', router)

app.listen(8080)