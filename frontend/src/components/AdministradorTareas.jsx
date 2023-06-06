import React, { useState } from 'react';
import Typography from '@mui/material/Typography';

function AdministradorTareas({AllProcesos,AllGenerales}) {
  // const [AllProcesos, setAllProcesos] = useState([
  //   {
  //     "Pid": 1,
  //     "Nombre": "systemd",
  //     "Usuario": 0,
  //     "Estado": 1,
  //     "Ram": 12,
  //     "Padre": 0
  //   },
  //   {
  //     "Pid": 2,
  //     "Nombre": "kthreadd",
  //     "Usuario": 0,
  //     "Estado": 1,
  //     "Ram": 0,
  //     "Padre": 1
  //   },
  //   {
  //     "Pid": 3,
  //     "Nombre": "rcu_gp",
  //     "Usuario": 0,
  //     "Estado": 1026,
  //     "Ram": 0,
  //     "Padre": 2
  //   },
  //   {
  //     "Pid": 4,
  //     "Nombre": "rcu_gp",
  //     "Usuario": 0,
  //     "Estado": 1026,
  //     "Ram": 0,
  //     "Padre": 0
  //   },
  // ]);

  const [procesosDesplegados, setProcesosDesplegados] = useState({});

  const handleDesplegarProcesos = (pid) => {
    setProcesosDesplegados((prev) => ({ ...prev, [pid]: !prev[pid] }));
  };

  const handleKill = async(pid) => {
    console.log(pid);
    try {
      const response = await fetch(`http://${process.env.REACT_APP_PUERTO}:8080/tasks`, {
        method: 'POST',
        body: pid, // Cuerpo de la solicitud POST, asegúrate de que sea un número entero válido
        headers: {
          'Content-Type': 'text/plain'
        }
      });
  
      if (response.ok) {
        console.log('La solicitud POST fue exitosa');
        // Realizar acciones adicionales si la solicitud es exitosa
      } else {
        console.log('La solicitud POST falló');
        // Realizar acciones adicionales si la solicitud falla
      }
    } catch (error) {
      console.log('Error al realizar la solicitud POST:', error);
      // Realizar acciones adicionales en caso de error
    }
  };


  const crearFilas = (procesos) => {
    return procesos.map((proceso) => (
      <React.Fragment key={proceso.pid}>
        <tr className={procesosDesplegados[proceso.pid] && 'desplegado'}>
          <td style={{backgroundColor: `${proceso.pid === 0 ||!AllProcesos.some((q) => q.pid === proceso.padre) ? "#FFE659": "" }`}} >{proceso.pid}</td>
          <td>{proceso.nombre}</td>
          <td>{proceso.usuario}</td>
          <td>{proceso.estado}</td>
          <td>{(proceso.ram/(AllGenerales[AllGenerales.length-1].totalram-Math.floor(Math.random() * (600 - 500 + 1) + 500))*100).toFixed(2)}</td>
          <td>{proceso.padre}</td>
          <td>
            {AllProcesos.filter((p) => p.padre === proceso.pid).length > 0 && (
              <button onClick={() => handleDesplegarProcesos(proceso.pid)}>
                {procesosDesplegados[proceso.pid] ? '-' : '+'}
              </button>
            )}
          </td>
          <td>
              <button onClick={() => handleKill(proceso.pid)}>
               x 
              </button>
          
          </td>
        </tr>
        {procesosDesplegados[proceso.pid] &&
          crearFilas(AllProcesos.filter((p) => p.padre === proceso.pid))}
      </React.Fragment>
    ));
  };
  

  return (
    <div>
      <table>
        <thead>
          <tr>
            <th style={{width: 200}} ><Typography variant="h5" color="inherit" component="div"><b>Pid</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Nombre</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Usuario</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Estado</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Ram (%)</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Padre</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Show More</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Kill</b></Typography></th>
          </tr>
        </thead>
        <tbody>
          {crearFilas(AllProcesos.filter((p) => p.padre === 0 || !AllProcesos.some((q) => q.pid === p.padre)) )}
        </tbody>
      </table>

    </div>
  );
}

export default AdministradorTareas;
