import axios from 'axios'

const baseURL = "/products/gloves/"

const getAll = async () => {
    console.log(baseURL)
    const res = await axios.get(baseURL)
    if (res.data === null)
        res.data = []
    return res.data
}

const accessoriesService = {
    getAll
}

export default accessoriesService
