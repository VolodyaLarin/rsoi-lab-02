import axios from "axios";


const timeout = 1000;
const failsCount = 5;

export function CreateBreakerContext(h) {
    const context = {
        _isOpen: true, _number: 0, _healthCheck: () => HealthTemplate(h), _url: h, _lastCall: null, _lastSuccessCall: null,
    }

    return context
}


function _retryTimeout(context) {
    setTimeout(async () => {
        const isH = await context._healthCheck()
        if (isH) {
            console.log(context._url, " is open")
            context._isOpen = true
            context._number = 0
        } else {
            _retryTimeout(context)
        }
    }, timeout)
}

export function DecorateBreaker(context, func) {
    console.log("Decorated ", context._url )
    return async (...a) => {
        if (!context._isOpen) {
            throw 'Service closed'
        }

        try {
            const res = await func(...a)
            context._number = 0
            context._lastCall = Date.now()
            context._lastSuccessCall = Date.now()
            return res
        } catch (err) {
            context._lastCall = Date.now()
            context._number++

            console.log(context._url, " has ", context._number, " fails")
            if (context._number >= failsCount) {
                console.log(context._url, " Closed")

                context._isOpen = false
                _retryTimeout(context)
            }

            throw err
        }
    }
}

async function HealthTemplate(url) {
    try {
        const res = await axios.get(url)
        return res.status === 200;
    } catch (err) {
        return false
    }
}
