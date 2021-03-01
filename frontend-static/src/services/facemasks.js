import axios from 'axios'

const baseURL = "/products/facemasks/"

const getAll = async () => {
    console.log(baseURL)
    const res = await axios.get(baseURL)
    if (res.data === null)
        res.data = []
    return res.data
}

const jacketsService = {
    getAll
}

export default jacketsService
