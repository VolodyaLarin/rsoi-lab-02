import axios from "axios";


const timeout = 5000;
const failsCount = 5;

export function CreateBreakerContext(h) {
    return {
        _isOpen: true, _number: 0, _lastCheck: Date.now(), _healthCheck: () => HealthTemplate(h)

    }
}


function _retryTimeout(context) {
    setTimeout(async () => {
        const isH = await context._healthCheck()
        if (isH) {
            context._isOpen = true
            context._number = 0
        } else {
            _retryTimeout(context)
        }
    }, timeout)
}

export function DecorateBreaker(context, func) {
    return (...a) => {
        if (!context._isOpen) {
            throw 'Service closed'
        }

        try {
            return func(...a)
        } catch (err) {
            context._number++
            if (context._number >= failsCount) {
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
