import React from 'react'
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { Line } from '@ant-design/plots';

   let dataTemp = [];
const Monitorram=({AllGenerales})=>{ 
   AllGenerales = AllGenerales.length > 0 ? AllGenerales : [{ramocupada:0,totalram:0}];
   dataTemp.push({ramocupada:AllGenerales[0].ramocupada,totalram:AllGenerales[0].totalram});

   const data = dataTemp.map((item,index) => {
      return {
        year: `T${index}`,
        value: (item.ramocupada/item.totalram*100 ).toFixed(2)
      }
   }).slice(-20);
  
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
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Monitor RAM
            </Typography>
            </center>
         </Paper>

         <br/>
         <br/>
         <div className='centerRam' >
            % Utilización RAM<br/><br/>
            <Line {...config}  />
         </div>


         <br/>
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Utilización
            </Typography>
            </center>
         </Paper>

         <br/>
         <center>
         <Typography variant="h5" color="inherit" component="div">
            {((AllGenerales[AllGenerales.length-1].ramocupada)/(AllGenerales[AllGenerales.length-1].totalram)*100).toFixed(2) } %
         </Typography>
         </center>

      </>
   );
}

export default Monitorram ;
