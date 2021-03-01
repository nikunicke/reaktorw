import React, { useState, useEffect } from 'react'
import glovesService from './services/gloves'
import beaniesService from './services/facemasks'
import facemasksService from './services/beanies'
import Container from 'react-bootstrap/Container'
import Tabs from 'react-bootstrap/Tabs'
import Tab from 'react-bootstrap/Tab'
import BootstrapTable from 'react-bootstrap-table-next'
import paginationFactory from 'react-bootstrap-table2-paginator'


const App = () => {
  const [ gloves, setGloves ] = useState([])
  const [ beanies, setBeanies ] = useState([])
  const [ facemasks, setFacemasks ] = useState([])

  useEffect(() => {
    glovesService.getAll().then(
      gloves => setGloves(gloves)
    )
  }, [])
  useEffect(() => {
    facemasksService.getAll().then(
      beanies => setBeanies(beanies)
    )
  }, [])
  useEffect(() => {
    beaniesService.getAll().then(
      facemasks => setFacemasks(facemasks)
    )
  }, [])



  const columns = [{
    dataField: "api_id",
    text: "ID"
  }, {
    dataField: "name",
    text: "Name"
  }, {
    dataField: "colors",
    text: "Colors",
    formatter: (cell) => {
      return cell.join(', ')
    }
  }, {
    dataField: "price",
    text: "Price"
  }, {
    dataField: "manufacturer",
    text: "Manufacturer"
  }, {
    dataField: "availability",
    text: "Availability",
    formatter: (cell) => (cell === "") ? "Please refresh" : cell ,
    style: function callback(cell, row, rowIndex, colIndex) {
      if (cell === "INSTOCK") {
        return ({ backgroundColor: "#5cb85c" })
      } else if (cell === "LESSTHAN10") {
        return ({ backgroundColor: "#f0ad4e" })
      } else if (cell === "OUTOFSTOCK") {
        return ({ backgroundColor: "#d9534f"})
      }
    }
  }]

  return (
    <div className="App">
      <Container fluid style={{maxWidth: "1300px"}}>
        <center><h1>reaktor warehouse</h1></center>
        <Tabs defaultActiveKey="gloves" id="uncontrolled-tab">
          <Tab eventKey="gloves" title="Gloves">
            <BootstrapTable
              bootstrap4 = {true}
              hover={true}
              striped={true}
              noDataIndication={() => "No data currently available. Try again shortly by pressing ctrl+R"}
              keyField="id"
              data={gloves}
              columns={columns}
              pagination={paginationFactory({})}/>
          </Tab>
          <Tab eventKey="beanies" title="Beanies">
           <BootstrapTable
              bootstrap4 = {true}
              hover={true}
              striped={true}
              noDataIndication={() => "No data currently available. Try again shortly by pressing ctrl+R"}
              keyField="id"
              data={beanies}
              columns={columns}
              pagination={paginationFactory({})}/>
          </Tab>
          <Tab eventKey="facemasks" title="Facemasks">
           <BootstrapTable
              bootstrap4 = {true}
              hover={true}
              striped={true}
              noDataIndication={() => "No data currently available. Try again shortly by pressing ctrl+R"}
              keyField="id"
              data={facemasks}
              columns={columns}
              pagination={paginationFactory({})}/>
          </Tab>
        </Tabs>
      </Container>
    </div>
  )
}

export default App;
