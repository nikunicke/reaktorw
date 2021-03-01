import axios from 'axios'

const baseURL = "/products/beanies/"

const getAll = async () => {
    console.log(baseURL)
    const res = await axios.get(baseURL)
    if (res.data === null)
        res.data = []
    return res.data
}

const shirtsService = {
    getAll
}

export default shirtsService
