#compdef _alicloud-vault alicloud-vault

# based on and inspirted by https://github.com/99designs/aws-vault/blob/master/contrib/completions/zsh/aws-vault.zsh

_alicloud-vault() {
    local i
    for (( i=2; i < CURRENT; i++ )); do
        if [[ ${words[i]} == -- ]]; then
            shift $i words
            (( CURRENT -= i ))
            _normal
            return
        fi
    done

    local matches=($(${words[1]} --completion-bash ${(@)words[2,$CURRENT]}))
    compadd -a matches

    if [[ $compstate[nmatches] -eq 0 && $words[$CURRENT] != -* ]]; then
        _files
    fi
}
