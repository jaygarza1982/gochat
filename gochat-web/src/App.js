import { Route, Switch, withRouter } from 'react-router-dom';
import './App.css';
import Chat from './components/Chat';

function App() {
  return (
    <div className="App">
      <Switch>
        <Route path='/chat' component={Chat} />
      </Switch>
    </div>
  );
}

export default withRouter(App);
