import React, { useEffect } from 'react'
import { makeStyles, withStyles } from '@material-ui/core/styles';
import BottomNavigation from '@material-ui/core/BottomNavigation';
import BottomNavigationAction from '@material-ui/core/BottomNavigationAction';
import Box from '@material-ui/core/Box';

import Game from './Game';

const useStyles = makeStyles({
  root: {
    width: '100%',
    height: 50,
    position: 'fixed',
    bottom: 0,
    background:"#c8c8c8"
  },
});

const Main = () => {

  const classes = useStyles();
  const [page, setPage] = React.useState(0);

  useEffect(
    () => {
    },[]
  )

  return (
    <React.Fragment>
      <Game/>
    </React.Fragment>
  );
}

export default Main;
