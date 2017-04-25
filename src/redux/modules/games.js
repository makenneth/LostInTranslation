import axios from 'axios';
import gameTemplate from 'helpers/gameTemplate';

const LOAD = 'hangperson/games/LOAD';
const LOAD_SUCCESS = 'hangperson/games/LOAD_SUCCESS';
const LOAD_FAIL = 'hangperson/games/LOAD_FAIL';
const CREATE = 'hangperson/games/CREATE';
const CREATE_SUCCESS = 'hangperson/games/CREATE_SUCCESS';
const CREATE_FAIL = 'hangperson/games/CREATE_FAIL';
const OTHER_CREATED = 'hangperson/games/OTHER_CREATED';
const CREATE_PRACTICE_GAME = 'hangperson/games/CREATE_PRACTICE_GAME';
const LOAD_PRACTICE_GAMES = 'hangperson/games/LOAD_PRACTICE_GAMES';

const initialState = {
  unfinished: [],
  finished: [],
  practice: [],
  loaded: false,
};
export default (state = initialState, action) => {
  switch (action.type) {
    case LOAD_SUCCESS:
      return {
        ...state,
        ...action.result.data,
        loaded: true,
      };
    case CREATE_SUCCESS:
      return {
        ...state,
        unfinished: [...state.unfinished, action.result.data],
      };
    case OTHER_CREATED:
      return {
        ...state,
        unfinished: [...state.unfinished, action.payload],
      };
    case CREATE_PRACTICE_GAME: {
      const practiceGameTemplate = gameTemplate(
        action.payload.id,
        action.payload.word,
        action.payload.user
      );
      console.log(practiceGameTemplate)
      return {
        ...state,
        practice: [
          ...state.practice,
          practiceGameTemplate,
        ],
      };
    }
    case LOAD_PRACTICE_GAMES:
      return {
        ...state,
        practice: action.payload,
      };
    default:
      return state;
  }
};

export const loadGames = () => {
  return {
    types: [LOAD, LOAD_SUCCESS, LOAD_FAIL],
    promise: axios.get('/api/user/games'),
  };
};

const createPracticeGameSuccess = (id, word, user) => {
  return {
    type: CREATE_PRACTICE_GAME,
    payload: {
      id,
      word,
      user,
    },
  };
}

export const createPracticeGame = (user) => {
  return (dispatch, getState) => {
    return axios.get('/api/games/random_word')
      .then(
        (res) => {
          const { practice } = getState().games;
          let maxId;
          if (practice.length > 0) {
            maxId = Math.max.apply(null, practice.map(g => g.info.id));
          } else {
            maxId = 0;
          }
          dispatch(createPracticeGameSuccess(maxId + 1, res.data.word, user));
          dispatch(savePracticeGame());
        },
        (err) => {
          console.warn(err);
        }
      );
  };
}

export const loadPracticeGames = () => {
  let practiceGames = [];
  const stored = localStorage.getItem('practiceGames');
  debugger;
  console.log(stored);
  if (stored) {
    practiceGames = JSON.parse(stored).filter(g => !g.info.finished);
  }
  return {
    type: LOAD_PRACTICE_GAMES,
    payload: practiceGames,
  };
};

export const createGame = (opponent) => {
  const req = axios({
    url: '/api/games/new',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json; charset=UTF-8',
    },
    data: {
      id: 0,
      userId2: opponent.id,
      username2: opponent.username,
      nickname2: opponent.nickname,
    },
  });
  return {
    types: [CREATE, CREATE_SUCCESS, CREATE_FAIL],
    promise: req,
  };
};

export const createdGame = (game) => {
  return {
    type: OTHER_CREATED,
    payload: game,
  };
};

export const isGamesLoaded = (state) => {
  return state.games.loaded;
};

export const savePracticeGame = () => {
  return (_, getState) => {
    const { games: { practice }, game, gameInfo } = getState();
    const currentGameIdx = practice.findIndex(g => g.info.id === gameInfo.id);
    let games = [];
    if (currentGameIdx > -1) {
      games = [
        ...practice.slice(0, currentGameIdx),
        { state: game, info: gameInfo },
        ...practice.slice(currentGameIdx + 1),
      ];
    } else {
      games = [ ...practice, { state: game, info: gameInfo } ];
    }
    localStorage.setItem('practiceGames', JSON.stringify(games));
  };
}

export const clearPracticeGame = (id) => {
  return (dispatch, getState) => {
    const { games: { practice } } = getState();
    if (practice.length) {
      const games = practice.filter(g => eval(g.info.id) !== eval(id));
      console.log(games);
      localStorage.setItem('practiceGames', JSON.stringify(games));
      dispatch(loadPracticeGames());
    }
  };
}
