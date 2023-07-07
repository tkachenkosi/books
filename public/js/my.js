
const KEY_ENTER = 'Esc';
const TIMEDELAY_DEFAULT = 500;

const elResult = document.querySelector('.result')
// const elAbtn = document.querySelector('.abtn')
const elLenResult = document.querySelector('#len_result')
const elHistoryResult = document.querySelector('#history_result')
const foundAvtor = document.querySelector('#foundAvtor')
const foundTitle = document.querySelector('#foundTitle')

const btnFoundAvtor = document.querySelector('#btnFoundAvtor')
const btnFoundTitle = document.querySelector('#btnFoundTitle')

let currPosHistoryInput = sessionStorage.length;

// установить занчение
/* const sq = function (e, s) {
    document.querySelector(e).innerHTML = s;
}; */

/* const sid = function (e, s) {
    document.getElementById(e).innerHTML = s;
}; */

const debounce = (cb, timeDelay = TIMEDELAY_DEFAULT) => {
  let timeoutId;
  return (...rest) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => cb.apply(this, rest), timeDelay);
  };
};


// находит книгу по номеру файла в архиве
const unZip = (e) => {
  // const code = e.textContent
  myAjax("/book", {code:e.textContent}, (data) => {
    if (data.book_autor === "err") {
      showAlert(data.book_name, 1)
    } else {
      showAlert("Книга "+e.textContent+".fb2 извлечена из архива "+data.fzip_name, 0)
    }
  });
}

// редактировать коментарий
/* const setRem = (id) => {
  showAlert("Книга "+id, 0)
  elResult.querySelector('#r'+id).insertAdjacentHTML('beforeend',`<input type="text" name="rem" value="" placeholder="коментарий">`)
} */

const hendlerHistory = () => {
  document.querySelectorAll('.abtn').forEach((e) => {
    e.removeEventListener('click', unZip);
  });
  elResult.innerHTML="" ;
  for (let i = 0; i < sessionStorage.length; i++) {
    const sesJSON = sessionStorage.getItem(sessionStorage.key(i))

      if (sesJSON === null) {
        return undefined
      }
      try {
        store = JSON.parse(sesJSON)
      } catch (e) {
        sessionStorage.removeItem(sessionStorage.key(i))
      }
      elResult.insertAdjacentHTML('beforeend', `<p>${store.autor} : ${store.name}</p>`)
  };
}


const hendlerFoundAvtor = () => {
  // console.log('Произошло событие', event.type)
  myAjax("/found", {autor:foundAvtor.value, name:foundTitle.value}, (data) => {

      elLenResult.innerHTML=data.length ;
      // предварительно удалить все ссылки что бы не было утечки памяти
      document.querySelectorAll('.abtn').forEach((e) => {
        e.removeEventListener('click', unZip);
      });

      elResult.innerHTML="" ;
      for (let i = 0; i < data.length; ++i) {
          // console.log("books ->", data[i].book_name, data[i].book_autor);
        elResult.insertAdjacentHTML('beforeend', `<p>${data[i].book_autor}<br><strong>${data[i].book_name}</strong> ${data[i].book_name2}<br><span class="section-info">${fDate(data[i].book_date)} ${data[i].book_lang} <a class="abtn" onclick="unZip(this)" href="#">${data[i].code}</a> ${data[i].fzip_name}</span></p>`);
        // elResult.insertAdjacentHTML('beforeend', `<p>${data[i].book_autor}<br><strong>${data[i].book_name}</strong> ${data[i].book_name2}<br><span class="section-info">${fDate(data[i].book_date)} ${data[i].book_lang} <a class="abtn" onclick="unZip(this)" href="#">${data[i].code}</a> ${data[i].fzip_name} <a class="abtn" onclick="setRem(${data[i].book_id})" href="#">rem</a></span><span id="r${data[i].book_id}" class="section-info"></span></p>`);

      };
      // сохраним параметры запроса
      sessionStorage.setItem(sessionStorage.length, JSON.stringify({autor:foundAvtor.value, name:foundTitle.value}))

      // const el = document.createElement('p');
      // el.textContent = data[1].book_name;
      // elResult.appendChild(el)
      // elResult.append(el)
      // elResult.insertAdjacentHTML('beforeend', '<p>'+data[2].book_name+'</p>');
    
  });
};

// возврат json
btnFoundAvtor.addEventListener('click', debounce(hendlerFoundAvtor));


const hendlerFoundTitle = () => {
  // console.log('Произошло событие', event.type)
  myAjaxHtml("/result", {autor:foundAvtor.value, name:foundTitle.value}, (data) => {
      // console.log("->", data);
      // elLenResult.innerHTML=data.length ;

      // предварительно удалить все ссылки что бы не было утечки памяти
      document.querySelectorAll('.abtn').forEach((e) => {
        e.removeEventListener('click', unZip);
      });
      sessionStorage.setItem(sessionStorage.length, JSON.stringify({autor:foundAvtor.value, name:foundTitle.value}))

      elResult.innerHTML=data;
  });

  setTimeout(setCountQuery, 2000);

}

// возрат шаблоном
btnFoundTitle.addEventListener('click', debounce(hendlerFoundTitle));

// количество книг в запросе (запускается с паузой)
function setCountQuery() {
  myAjaxHtml("/count", {}, (data) => {
    elLenResult.innerHTML=data;
  });
};


btnClearAvtor.addEventListener('click', (event) => {
  foundAvtor.value = "";
});

btnMoveAvtor.addEventListener('click', (event) => {
  currPosHistoryInput = (currPosHistoryInput <= 0) ? 0 : --currPosHistoryInput;

  if (currPosHistoryInput >= 0 && currPosHistoryInput <= sessionStorage.length) {
    const storeJSON = sessionStorage.getItem(currPosHistoryInput)

    if (storeJSON === null) {
      return undefined
    }
    try {
      store = JSON.parse(storeJSON)
    } catch (e) {
      sessionStorage.removeItem(currPosHistoryInput)
      return undefined
    }

    foundAvtor.value = store.autor;
    foundTitle.value = store.name;
  }
  elHistoryResult.innerHTML=`${currPosHistoryInput} / ${sessionStorage.length}`;

});

btnClearTitle.addEventListener('click', (event) => {
  foundTitle.value = "";
});

btnMoveTitle.addEventListener('click', (event) => {
/*
  currPosHistoryInput = (currPosHistoryInput >= sessionStorage.length) ? sessionStorage.length : ++currPosHistoryInput;

  if (currPosHistoryInput >= 0 && currPosHistoryInput <= sessionStorage.length) {
    const storeJSON = sessionStorage.getItem(currPosHistoryInput)

    if (storeJSON === null) {
      return undefined
    }
    try {
      store = JSON.parse(storeJSON)
    } catch (e) {
      sessionStorage.removeItem(currPosHistoryInput)
      return undefined
    }

    foundAvtor.value = store.autor;
    foundTitle.value = store.name;
  }
  // показать позицию в истории поиска
  elHistoryResult.innerHTML=`${currPosHistoryInput} /  ${sessionStorage.length}`;
 */

  // показать историю
  hendlerHistory()
});


foundAvtor.addEventListener('keyup', (event) => {
  if (event.key === "Enter") {
    debounce(hendlerFoundAvtor())
    // showAlert("Проверка работы сообщения")
  }
});

foundTitle.addEventListener('keyup', (event) => {
  if (event.key === "Enter") {
    debounce(hendlerFoundAvtor())
    // showAlert("Проверка работы сообщения")
  }
});

// мой универсальный ajax возврат json
function myAjax(url, b, cback) {
    fetch(url, {
        method: 'post',
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        },
        body: JSON.stringify(b)
    }).then((response) => {
        if (response.status !== 200) {
            console.log('Проблема! Status Code: ' + response.status);
            return;
        };
        return response.json(); // конвертируем ответ в json
    }).then((data) => {
        cback(data); // а вот и наши данные в формате json
    }).catch((error) => {
        console.log('Request myAjax ошибка:', error);
    });
};

// мой универсальный ajax возврат html
// headers: {"Content-type": "text/html; charset=UTF-8"}, 
function myAjaxHtml(url, b, cback) {
    fetch(url, {
        method: 'post',
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        },
        body: JSON.stringify(b)
    }).then((response) => {
        if (response.status !== 200) {
            console.log('Проблема! Status Code: ' + response.status);
            return;
        };
        return response.text() // возврать в формате html
    }).then((data) => {
        cback(data);
    }).catch((error) => {
        console.log('Request ошибка:', error);
    });
};

function fDate(dat) {
  dt = new Date(dat);
  const y=dt.getFullYear();
  let m=''+(dt.getMonth()+1);
  let d=''+dt.getDate();

  if (d.length<2) d='0'+d;
  if (m.length<2) m='0'+m;

  return [y, m, d].join('-');
};


const isEnter = (evt) => evt.key === "Enter";

const showAlert = (message, type) => {
  const SHOW_ALERT_TIME = 3000;
  const alertContainer = document.createElement('div');
  alertContainer.style.zIndex = '1000';
  alertContainer.style.position = 'fixed';
  alertContainer.style.left = '0';
  alertContainer.style.top = '0';
  alertContainer.style.right = '0';
  alertContainer.style.padding = '10px 3px';
  alertContainer.style.fontSize = '21px';
  alertContainer.style.color = '#c4d0e8';
  alertContainer.style.textAlign = 'center';
  if (type === 1) {
    alertContainer.style.backgroundColor = '#aa3104';
  } else {
    alertContainer.style.backgroundColor = '#4b70b5';
  }

  alertContainer.textContent = message;

  document.body.append(alertContainer);

  setTimeout(() => {
    alertContainer.remove();
  }, SHOW_ALERT_TIME);
};
