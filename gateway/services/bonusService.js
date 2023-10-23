import axios from "axios";
import {genUsernameHeaders, axiosLogger, axiosLoggerResp} from "./utils.js";
import {CreateBreakerContext, DecorateBreaker} from "./cicuitBreakerDecorator.js";
import * as fs from "fs";

const timeout = 2000;


const baseURL = process.env.BONUS_SERVICE_URL || "http://localhost:8050/api/v1/bonus";

const BonusService = {
    _tasks: [],

    _service: new axios.Axios({
        baseURL
    }), GetBonusDetails: async function (username) {
        const request = await this._service.get(`/`, genUsernameHeaders(username));
        if (request.status !== 200) {
            throw "Request error"
        }
        return JSON.parse(request.data)
    }, CreateBonusItem: async function (flightNumber, price, paidFromBalance, ticketUid, username) {
        const request = await this._service.post('/', JSON.stringify({
            flightNumber, price, paidFromBalance, ticketUid: ticketUid
        }), genUsernameHeaders(username));

        return JSON.parse(request.data)
    }, _deleteBonusItem: async function (uid, username) {
        await this._service.delete(`/${uid}`, genUsernameHeaders(username));
    }, DeleteBonusItem: async function (uid, username) {
        try {
            await this._deleteBonusItem(uid, username)
        } catch (err) {
            this._addTask([uid, username])
        }

        return
    },
    _addTask(task) {
        const details = {
            task,
            id: Date.now() % 1000000 + Math.random()
        }
        this._tasks.push(details)
        console.log("New task ", details)
    }
}


BonusService._service.interceptors.request.use(axiosLogger)
BonusService._service.interceptors.response.use(axiosLoggerResp)

const ctx = CreateBreakerContext((new URL(baseURL).origin + '/manage/health'))

BonusService.GetBonusDetails = DecorateBreaker(ctx, BonusService.GetBonusDetails.bind(BonusService))
BonusService.CreateBonusItem = DecorateBreaker(ctx, BonusService.CreateBonusItem.bind(BonusService))
BonusService._deleteBonusItem = DecorateBreaker(ctx, BonusService._deleteBonusItem.bind(BonusService))

// setInterval(() => {
//     console.log(ctx)
// }, 2000)

setInterval(async () => {
    const tasksCopy = [...BonusService._tasks]
    // console.log("Process ", tasksCopy)

    for (let i in tasksCopy) {
        const task = tasksCopy[i]

        console.log("Do bonuses task", task)
        try {
            await BonusService.DeleteBonusItem(...task.task)
            BonusService._tasks = BonusService._tasks.filter(x => x.id !== task.id)
        } catch (err) {
            console.log("Error bonuses task", task)
        }

    }


}, timeout)


export {BonusService}
