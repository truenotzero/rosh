const note_selected = (e) => {
    const note = document.getElementById("active-note");
    note.innerHTML = e;
    note.value = e;

    const save = document.getElementById("save");
    save.hidden = false;

    const text_input = document.getElementById("text-input");
    text_input.disabled = false;
};

let is_dirty = false;
const kbhit = () => {
    is_dirty = true
    document.getElementById("dirty").hidden = false;
};

const save = () => {
    is_dirty = false
    document.getElementById("dirty").hidden = true;
};

document.addEventListener("DOMContentLoaded", (event) => {
    document.body.addEventListener("htmx:beforeSwap", function (evt) {
        const elt = evt.detail.elt;
        if (elt.id === "text-input" && evt.detail.xhr.status == 200) {
            evt.detail.shouldSwap = false;
            elt.value = evt.detail.serverResponse;
        }
    });

    document.body.addEventListener("htmx:confirm", function (evt) {
        if (!is_dirty) {
            evt.preventDefault();
            evt.detail.issueRequest(true);
        }

        save();
    });
});

const active_note = () => {
    const note = document.getElementById("active-note");
    return note.innerHTML;
};

const text_input = () => {
    const inp = document.getElementById("text-input");
    return inp.innerHTML;
};

const toggle_word_wrap = () => {
    const ipt = document.getElementById("text-input");
    if (ipt.wrap === "off") {
        ipt.wrap = "soft";
    } else {
        ipt.wrap = "off";
    }
}

const toggle_spell_check = (checkbox) => {
    const ipt = document.getElementById("text-input");
    ipt.spellcheck = checkbox.checked;
}

const toggle_rtl = (ckb) => {
    const ipt = document.getElementById("text-input");
    ipt.dir = ckb.checked ? "rtl" : "ltr";
}
