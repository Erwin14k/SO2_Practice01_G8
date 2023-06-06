import React from 'react'
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { Line } from '@ant-design/plots';
import {CanvasJSChart} from 'canvasjs-react-charts'

import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';


let dataTemp = [];
const Monitorram = ({ AllGenerales }) => {
   AllGenerales = AllGenerales.length > 0 ? AllGenerales : [{ ramocupada: 0, totalram: 0 }];
   dataTemp.push({ ramocupada: AllGenerales[0].ramocupada, totalram: AllGenerales[0].totalram });

   const data = dataTemp.map((item, index) => {
      return {
         year: `T${index}`,
         value: (item.ramocupada).toFixed(2)
      }
   }).slice(-20);

   const data2 = dataTemp.map((item, index) => {
      return {
         y: `T${index}`,
         x: (item.ramocupada).toFixed(2)
      }
   }).slice(-20);
   
   const options = {
      animationEnabled: true,
      exportEnabled: true,
      theme: "light2", // "light1", "dark1", "dark2"
      title:{
         text: "Bounce Rate by Week of Year"
      },
      axisY: {
         title: "Bounce Rate",
         suffix: "%"
      },
      axisX: {
         title: "Week of Year",
         prefix: "W",
         interval: 2
      },
      data: [{
         type: "line",
         toolTipContent: "Week {x}: {y}%",
         dataPoints: data2
      }]
   }



   const config = {
      data,
      autoFit: false,
      xField: 'year',
      yField: 'value',
      point: {
         size: 5,
         shape: 'diamond',
      },
      label: {
         style: {
            fill: '#aaa',
         },
      },
      color: '#D921F7'
   };



   return (
      <>
         <br />
         <Paper >
            <center>
               <Typography variant="h4" color="inherit" component="div">
                  Monitor RAM
               </Typography>
            </center>
         </Paper>

         <br />
         <br />
         <div>
        
      </div>
      <CanvasJSChart options = {options}/>
         {/* <div className='centerRam' >
             Utilización RAM (MB)<br /><br />
   <Line {...config} /> 
    </div>*/}
        

       
    
         <br />
         <br />
         <Paper >
            <center>
               <Typography variant="h4" color="inherit" component="div">
                  RAM
               </Typography>
            </center>
         </Paper>

         <br />
         <br />

         <Table aria-label="simple table" >
            <TableHead style={{ backgroundColor: "#FFE659" }}>
               <TableRow >
                  <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>%Usado</b></Typography></TableCell>
                  <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Total </b></Typography></TableCell>
                  <TableCell align="center"><Typography variant="h5" color="inherit" component="div"><b>Consumida </b></Typography></TableCell>
               </TableRow>
            </TableHead>
            <TableBody>
               <TableRow >
                  <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{((AllGenerales[AllGenerales.length - 1].ramocupada) / (AllGenerales[AllGenerales.length - 1].totalram) * 100).toFixed(2)} %</Typography></TableCell>
                  <TableCell align="center"><Typography variant="h6" color="inherit" component="div">{((AllGenerales[AllGenerales.length - 1].totalram)).toFixed(2)} MB</Typography></TableCell>
                  <TableCell align="center"><Typography variant="h6" color="inherit" component="div"> {((AllGenerales[AllGenerales.length - 1].ramocupada)).toFixed(2)} MB</Typography></TableCell>
               </TableRow>
            </TableBody>
         </Table>


         <br />
         <br />


      </>
   );
}

export default Monitorram;
