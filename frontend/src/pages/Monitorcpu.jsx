import React from 'react'
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { Line } from '@ant-design/plots';


   let dataTemp = [];

const Monitorcpu=({AllGenerales})=>{ 
   AllGenerales = AllGenerales.length > 0 ? AllGenerales : [{totalcpu:0}];
   dataTemp.push({totalcpu:AllGenerales[0].totalcpu});


   const data = dataTemp.map((item,index) => {
      return {
        year: `T${index}`,
        value: item.totalcpu
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
      color: '#4AC86F'
    };
  

   return (
      <>
         <br/>
         <Paper >
            <center>
            <Typography variant="h4" color="inherit" component="div">
               Monitor CPU
            </Typography>
            </center>
         </Paper>

         <br/>
         <br/>
         <div className='center' >
            % Utilización CPU<br/><br/>
            <Line {...config} />
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
            {AllGenerales[AllGenerales.length-1].totalcpu} %
         </Typography>
         </center>

      </>
   );
}

export default Monitorcpu ;
