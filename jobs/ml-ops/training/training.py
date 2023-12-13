import pandas as pd
import numpy as np
from imblearn.over_sampling import SMOTE
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score, precision_score, recall_score, log_loss
from sklearn.model_selection import RandomizedSearchCV
from sklearn.ensemble import RandomForestClassifier


def transform_data(data: pd.DataFrame) -> pd.DataFrame:
    """Handles the transformation of categorical variables of the dataset into 0/1 indicators"""

    # Use the same category for basic education sub-categories
    data["education"] = np.where(
        data["education"] == "basic.9y", "Basic", data["education"]
    )
    data["education"] = np.where(
        data["education"] == "basic.6y", "Basic", data["education"]
    )
    data["education"] = np.where(
        data["education"] == "basic.4y", "Basic", data["education"]
    )

    # Transform categorical variables into 0/1 indicators and remove columns with string categories
    cat_vars = [
        "job",
        "marital",
        "education",
        "default",
        "housing",
        "loan",
        "contact",
        "month",
        "day_of_week",
        "poutcome",
    ]
    for var in cat_vars:
        cat_list = "var" + "_" + var
        cat_list = pd.get_dummies(data[var], prefix=var)
        data = data.join(cat_list)

    data_vars = data.columns.values.tolist()
    to_keep = [i for i in data_vars if i not in cat_vars]
    data = data[to_keep]

    # Normalize column naming
    data.columns = data.columns.str.replace(".", "_")
    data.columns = data.columns.str.replace(" ", "_")

    # Replace yes/no by 1/0 for target variable y
    data["y"] = data["y"].replace(to_replace=["yes", "no"], value=[1, 0])

    return data


def split_to_train_test_data(
    data: pd.DataFrame,
) -> tuple[pd.DataFrame, pd.DataFrame, pd.DataFrame, pd.DataFrame]:
    """Extracts target (predicted) variable and splits data into training/test data"""

    # Extract target (predicted) variable
    x = data.loc[:, data.columns != "y"]
    y = data.loc[:, data.columns == "y"]

    x_train, x_test, y_train, y_test = train_test_split(
        x, y, test_size=0.2, stratify=y, random_state=50
    )

    return x_train, x_test, y_train, y_test


def over_sample_target_class(
    x_train: pd.DataFrame, y_train: pd.DataFrame
) -> tuple[pd.DataFrame, pd.DataFrame]:
    """Resamples training data"""

    over_sampler = SMOTE(random_state=0)
    sampled_data_x, sampled_data_y = over_sampler.fit_resample(x_train, y_train)

    return pd.DataFrame(data=sampled_data_x, columns=x_train.columns), pd.DataFrame(
        data=sampled_data_y, columns=["y"]
    )


def compute_performance_metrics(
    y_true: pd.DataFrame, y_pred: np.ndarray, y_pred_prob: np.ndarray
) -> dict:
    """Computes performance metrics"""

    acc = accuracy_score(y_true, y_pred)
    prec = precision_score(y_true, y_pred)
    recall = recall_score(y_true, y_pred)
    entropy = log_loss(y_true, y_pred_prob)

    return {
        "accuracy": round(acc, 2),
        "precision": round(prec, 2),
        "recall": round(recall, 2),
        "entropy": round(entropy, 2),
    }


def tune_classifier(
    x_train: pd.DataFrame, y_train: np.ndarray
) -> tuple[RandomForestClassifier, dict]:
    """Looks for optimal classifier hyperparameters then use them to fit a classifier"""

    random_grid = {
        "n_estimators": [5, 21, 51, 101],
        "max_features": ["sqrt"],
        "max_depth": [int(x) for x in np.linspace(10, 120, num=12)],
        "min_samples_split": [2, 6, 10],
        "min_samples_leaf": [1, 3, 4],
        "bootstrap": [True, False],
    }

    classifier = RandomForestClassifier()
    random_search = RandomizedSearchCV(
        estimator=classifier,
        param_distributions=random_grid,
        n_iter=100,
        cv=5,
        verbose=2,
        random_state=40,
        n_jobs=-1,
    )
    random_search.fit(x_train, y_train)

    optimal_params = random_search.best_params_
    n_estimators = optimal_params["n_estimators"]
    min_samples_split = optimal_params["min_samples_split"]
    min_samples_leaf = optimal_params["min_samples_leaf"]
    max_features = optimal_params["max_features"]
    max_depth = optimal_params["max_depth"]
    bootstrap = optimal_params["bootstrap"]

    opt_classifier = RandomForestClassifier(
        n_estimators=n_estimators,
        min_samples_split=min_samples_split,
        min_samples_leaf=min_samples_leaf,
        max_features=max_features,
        max_depth=max_depth,
        bootstrap=bootstrap,
    )
    opt_classifier.fit(x_train, y_train)

    return opt_classifier, optimal_params
